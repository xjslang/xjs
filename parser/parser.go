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

// Interceptor is a generic function type that allows middleware-style interception
// of parsing operations. It receives the parser and a next function to call the
// original parsing logic, enabling extension and modification of parsing behavior.
type Interceptor[T ast.Node] func(p *Parser, next func() T) T

// Transformer is a function type that takes an AST Program and returns a modified version.
// It's used for post-processing the entire parsed program.
type Transformer func(*ast.Program) *ast.Program

// Builder provides a fluent interface for constructing a Parser with various middleware
// and transformations. It allows composition of lexer options, statement interceptors,
// expression interceptors, and program transformers.
type Builder struct {
	LexerBuilder     *lexer.Builder
	stmtInterceptors []Interceptor[ast.Statement]
	expInterceptors  []Interceptor[ast.Expression]
}

// Parser is the main structure responsible for syntactic analysis of XJS source code.
// It implements a Pratt parser (top-down operator precedence parser) that converts
// a stream of tokens from the lexer into an Abstract Syntax Tree (AST).
//
// The parser supports middleware-style extensions through interceptors and transformers,
// allowing dynamic modification of parsing behavior without changing the core parser logic.
type Parser struct {
	// lexer provides the stream of tokens to be parsed
	lexer *lexer.Lexer

	// CurrentToken is the token currently being processed
	CurrentToken token.Token
	// PeekToken is the next token in the stream (lookahead)
	PeekToken token.Token

	// statementParseFn is the current statement parsing function (can be intercepted)
	statementParseFn func(*Parser) ast.Statement
	// expressionParseFn is the current expression parsing function (can be intercepted)
	expressionParseFn func(*Parser, int) ast.Expression
	// prefixParseFns maps token types to their prefix parsing functions
	prefixParseFns map[token.Type]func() ast.Expression
	// infixParseFns maps token types to their infix parsing functions
	infixParseFns map[token.Type]func(ast.Expression) ast.Expression

	// errors accumulates parsing errors encountered during parsing
	errors []ParserError

	// Context stack for tracking parsing state
	contextStack []ContextType
	// Current expression precedence during parsing
	currentExpressionPrecedence int
}

type parserOptions struct {
	stmtInterceptors []Interceptor[ast.Statement]
	expInterceptors  []Interceptor[ast.Expression]
}

// newWithOptions creates a new Parser instance with the specified lexer and parser options.
// It initializes all parsing function maps, sets up prefix and infix parsers for various
// token types, applies any provided interceptors and transformers, and prepares the parser
// for parsing by reading the first two tokens.
//
// Parameters:
//   - l: The lexer that will provide tokens for parsing
//   - opts: Configuration options including interceptors and transformers
//
// Returns a fully initialized Parser ready to parse source code.
func newWithOptions(l *lexer.Lexer, opts parserOptions) *Parser {
	p := &Parser{
		lexer:             l,
		errors:            []ParserError{},
		statementParseFn:  baseParseStatement,
		expressionParseFn: baseParseExpression,
		contextStack:      []ContextType{GlobalContext}, // Initialize with global context
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

	// Read two tokens, so CurrentToken and PeekToken are both set
	p.NextToken()
	p.NextToken()

	return p
}

// New creates a new Parser instance with default options.
// This is the standard constructor for creating a parser with no middleware or transformations.
func New(l *lexer.Lexer) *Parser {
	return newWithOptions(l, parserOptions{})
}

// ParseProgram parses the entire source code and returns the resulting AST Program.
// It continuously parses statements until EOF is reached, collecting any parsing errors.
// After parsing, it applies any registered program transformers to the final AST.
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
	if len(p.errors) > 0 {
		return program, fmt.Errorf("parsing failed with %d errors: %v",
			len(p.errors), p.errors[0])
	}
	return program, nil
}

// NextToken advances the parser to the next token in the stream.
// It moves PeekToken to CurrentToken and reads a new token from the lexer into PeekToken.
// This maintains the one-token lookahead that enables efficient parsing decisions.
func (p *Parser) NextToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

// AddError creates and adds a new parsing error to the parser's error collection.
// It captures the current token's position information and creates a structured
// ParserError with the provided message.
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

// ExpectToken checks if the next token (PeekToken) matches the expected type.
// If it matches, the parser advances to that token and returns true.
// If it doesn't match, it adds a parsing error and returns false.
func (p *Parser) ExpectToken(t token.Type) bool {
	if p.PeekToken.Type == t {
		p.NextToken()
		return true
	}
	p.AddError(fmt.Sprintf("output %s, got %s", t, p.PeekToken.Type))
	return false
}

// Errors returns a copy of all parsing errors encountered during parsing.
// This allows external code to inspect and handle parsing errors appropriately.
func (p *Parser) Errors() []ParserError {
	return p.errors
}

// peekPrecedence returns the operator precedence of the next token (PeekToken).
// This is used in the Pratt parser algorithm to determine whether to continue
// parsing an expression or return to a higher precedence level.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.PeekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// currentPrecedence returns the operator precedence of the current token.
// This is used in the Pratt parser algorithm for precedence comparison
// during expression parsing.
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.CurrentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// useStatementInterceptor applies a middleware interceptor to statement parsing.
// It wraps the current statement parsing function with the provided interceptor,
// enabling modification or extension of statement parsing behavior.
func (p *Parser) useStatementInterceptor(interceptor Interceptor[ast.Statement]) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return interceptor(p, func() ast.Statement {
			return next(p)
		})
	}
}

// useExpressionInterceptor applies a middleware interceptor to expression parsing.
// It wraps the current expression parsing function with the provided interceptor,
// enabling modification or extension of expression parsing behavior.
// It also manages the current expression precedence context during interception.
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
