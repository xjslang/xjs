// Package parser provides syntax analysis functionality for the XJS language.
// It builds an Abstract Syntax Tree (AST) from tokens provided by the lexer.
package parser

import (
	"fmt"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

// Operator precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGNMENT  // =
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALITY    // == !=
	COMPARISON  // > < >= <=
	SUM         // +
	PRODUCT     // * / %
	UNARY       // -x !x ++x --x
	POSTFIX     // x++ x--
	CALL        // myFunction(X)
	MEMBER      // obj.prop obj[prop]
)

var precedences = map[token.Type]int{
	token.ASSIGN:       ASSIGNMENT,
	token.PLUS_ASSIGN:  ASSIGNMENT,
	token.MINUS_ASSIGN: ASSIGNMENT,
	token.OR:           LOGICAL_OR,
	token.AND:          LOGICAL_AND,
	token.EQ:           EQUALITY,
	token.NOT_EQ:       EQUALITY,
	token.LT:           COMPARISON,
	token.GT:           COMPARISON,
	token.LTE:          COMPARISON,
	token.GTE:          COMPARISON,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.MULTIPLY:     PRODUCT,
	token.DIVIDE:       PRODUCT,
	token.MODULO:       PRODUCT,
	token.INCREMENT:    POSTFIX,
	token.DECREMENT:    POSTFIX,
	token.LPAREN:       CALL,
	token.DOT:          MEMBER,
	token.LBRACKET:     MEMBER,
}

type Interceptor[T ast.Node] func(p *Parser, next func() T) T
type Transformer func(*ast.Program) *ast.Program

type Builder struct {
	LexerBuilder        *lexer.Builder
	stmtInterceptors    []Interceptor[ast.Statement]
	expInterceptors     []Interceptor[ast.Expression]
	programTransformers []Transformer
}

type Parser struct {
	lexer *lexer.Lexer

	CurrentToken token.Token
	PeekToken    token.Token

	statementParseFn   func(*Parser) ast.Statement
	expressionParseFn  func(*Parser, int) ast.Expression
	transformProgramFn func(*ast.Program) *ast.Program
	prefixParseFns     map[token.Type]func() ast.Expression
	infixParseFns      map[token.Type]func(ast.Expression) ast.Expression

	errors []ParserError

	// Context stack for tracking parsing state
	contextStack []ContextType
	// Current expression precedence during parsing
	currentExpressionPrecedence int
}

type parserOptions struct {
	stmtInterceptors    []Interceptor[ast.Statement]
	expInterceptors     []Interceptor[ast.Expression]
	programTransformers []Transformer
}

func newWithOptions(l *lexer.Lexer, opts parserOptions) *Parser {
	p := &Parser{
		lexer:              l,
		errors:             []ParserError{},
		statementParseFn:   baseParseStatement,
		expressionParseFn:  baseParseExpression,
		transformProgramFn: nil,
		contextStack:       []ContextType{GlobalContext}, // Initialize with global context
	}

	p.prefixParseFns = make(map[token.Type]func() ast.Expression)
	p.prefixParseFns[token.IDENT] = p.ParseIdentifier
	p.prefixParseFns[token.INT] = p.ParseIntegerLiteral
	p.prefixParseFns[token.FLOAT] = p.ParseFloatLiteral
	p.prefixParseFns[token.STRING] = p.ParseStringLiteral
	p.prefixParseFns[token.RAW_STRING] = p.ParseMultiStringLiteral
	p.prefixParseFns[token.TRUE] = p.ParseBooleanLiteral
	p.prefixParseFns[token.FALSE] = p.ParseBooleanLiteral
	p.prefixParseFns[token.NULL] = p.ParseNullLiteral
	p.prefixParseFns[token.NOT] = p.ParseUnaryExpression
	p.prefixParseFns[token.MINUS] = p.ParseUnaryExpression
	p.prefixParseFns[token.INCREMENT] = p.ParseUnaryExpression
	p.prefixParseFns[token.DECREMENT] = p.ParseUnaryExpression
	p.prefixParseFns[token.LPAREN] = p.ParseGroupedExpression
	p.prefixParseFns[token.LBRACKET] = p.ParseArrayLiteral
	p.prefixParseFns[token.LBRACE] = p.ParseObjectLiteral
	p.prefixParseFns[token.FUNCTION] = p.ParseFunctionExpression

	p.infixParseFns = make(map[token.Type]func(ast.Expression) ast.Expression)
	p.infixParseFns[token.PLUS] = p.ParseBinaryExpression
	p.infixParseFns[token.MINUS] = p.ParseBinaryExpression
	p.infixParseFns[token.MULTIPLY] = p.ParseBinaryExpression
	p.infixParseFns[token.DIVIDE] = p.ParseBinaryExpression
	p.infixParseFns[token.MODULO] = p.ParseBinaryExpression
	p.infixParseFns[token.EQ] = p.ParseBinaryExpression
	p.infixParseFns[token.NOT_EQ] = p.ParseBinaryExpression
	p.infixParseFns[token.LT] = p.ParseBinaryExpression
	p.infixParseFns[token.GT] = p.ParseBinaryExpression
	p.infixParseFns[token.LTE] = p.ParseBinaryExpression
	p.infixParseFns[token.GTE] = p.ParseBinaryExpression
	p.infixParseFns[token.AND] = p.ParseBinaryExpression
	p.infixParseFns[token.OR] = p.ParseBinaryExpression
	p.infixParseFns[token.ASSIGN] = p.ParseAssignmentExpression
	p.infixParseFns[token.PLUS_ASSIGN] = p.ParseCompoundAssignmentExpression
	p.infixParseFns[token.MINUS_ASSIGN] = p.ParseCompoundAssignmentExpression
	p.infixParseFns[token.LPAREN] = p.ParseCallExpression
	p.infixParseFns[token.DOT] = p.ParseMemberExpression
	p.infixParseFns[token.LBRACKET] = p.ParseComputedMemberExpression
	p.infixParseFns[token.INCREMENT] = p.ParsePostfixExpression
	p.infixParseFns[token.DECREMENT] = p.ParsePostfixExpression

	// uses middlewares
	for _, interceptor := range opts.stmtInterceptors {
		p.useStatementInterceptor(interceptor)
	}
	for _, interceptor := range opts.expInterceptors {
		p.useExpressionInterceptor(interceptor)
	}
	for _, transformer := range opts.programTransformers {
		p.useProgramTransformer(transformer)
	}

	// Read two tokens, so CurrentToken and PeekToken are both set
	p.NextToken()
	p.NextToken()

	return p
}

func New(l *lexer.Lexer) *Parser {
	return newWithOptions(l, parserOptions{})
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.CurrentToken.Type != token.EOF {
		stmt := p.statementParseFn(p)
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	if p.transformProgramFn != nil {
		program = p.transformProgramFn(program)
	}
	if len(p.errors) > 0 {
		return program, fmt.Errorf("parsing failed with %d errors: %v",
			len(p.errors), p.errors[0])
	}
	return program, nil
}

func (p *Parser) NextToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

func (p *Parser) AddError(message string) {
	pos := Position{
		Line:   p.CurrentToken.Line,
		Column: p.CurrentToken.Column,
	}
	err := ParserError{
		Message:  message,
		Position: pos,
		Code:     "SYNTAX_ERROR",
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) ExpectToken(t token.Type) bool {
	if p.PeekToken.Type == t {
		p.NextToken()
		return true
	}
	p.AddError(fmt.Sprintf("output %s, got %s", t, p.PeekToken.Type))
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.PeekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.CurrentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) useStatementInterceptor(interceptor Interceptor[ast.Statement]) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return interceptor(p, func() ast.Statement {
			return next(p)
		})
	}
}

func (p *Parser) useExpressionInterceptor(interceptor Interceptor[ast.Expression]) {
	next := p.expressionParseFn
	p.expressionParseFn = func(p *Parser, precedence int) ast.Expression {
		oldPrecedence := p.currentExpressionPrecedence
		p.currentExpressionPrecedence = precedence
		defer func() {
			p.currentExpressionPrecedence = oldPrecedence
		}()
		return interceptor(p, func() ast.Expression {
			return next(p, precedence)
		})
	}
}

func (p *Parser) useProgramTransformer(transformer Transformer) {
	next := p.transformProgramFn
	p.transformProgramFn = func(program *ast.Program) *ast.Program {
		program = transformer(program)
		if next == nil {
			return program
		}
		return next(program)
	}
}
