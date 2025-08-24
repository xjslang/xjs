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

// Token precedences mapping
var precedences = map[token.Type]int{
	token.ASSIGN:        ASSIGNMENT,
	token.OR:            LOGICAL_OR,
	token.AND:           LOGICAL_AND,
	token.EQ:            EQUALITY,
	token.NOT_EQ:        EQUALITY,
	token.EQ_STRICT:     EQUALITY,
	token.NOT_EQ_STRICT: EQUALITY,
	token.LT:            COMPARISON,
	token.GT:            COMPARISON,
	token.LTE:           COMPARISON,
	token.GTE:           COMPARISON,
	token.PLUS:          SUM,
	token.MINUS:         SUM,
	token.MULTIPLY:      PRODUCT,
	token.DIVIDE:        PRODUCT,
	token.MODULO:        PRODUCT,
	token.LPAREN:        CALL,
	token.DOT:           MEMBER,
	token.LBRACKET:      MEMBER,
}

// Parser represents the parser state and configuration
type Parser struct {
	lexer *lexer.Lexer

	CurrentToken token.Token
	PeekToken    token.Token

	parseStatement func(p *Parser) ast.Statement

	prefixParseFns map[token.Type]func() ast.Expression
	infixParseFns  map[token.Type]func(ast.Expression) ast.Expression

	errors []string
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:          l,
		errors:         []string{},
		parseStatement: baseParseStatement,
	}

	// Initialize prefix parse functions
	p.prefixParseFns = make(map[token.Type]func() ast.Expression)
	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.FLOAT] = p.parseFloatLiteral
	p.prefixParseFns[token.STRING] = p.parseStringLiteral
	p.prefixParseFns[token.TRUE] = p.parseBooleanLiteral
	p.prefixParseFns[token.FALSE] = p.parseBooleanLiteral
	p.prefixParseFns[token.NULL] = p.parseNullLiteral
	p.prefixParseFns[token.NOT] = p.parseUnaryExpression
	p.prefixParseFns[token.MINUS] = p.parseUnaryExpression
	p.prefixParseFns[token.INCREMENT] = p.parseUnaryExpression
	p.prefixParseFns[token.DECREMENT] = p.parseUnaryExpression
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression
	p.prefixParseFns[token.LBRACKET] = p.parseArrayLiteral
	p.prefixParseFns[token.LBRACE] = p.parseObjectLiteral
	p.prefixParseFns[token.FUNCTION] = p.parseFunctionExpression

	// Initialize infix parse functions
	p.infixParseFns = make(map[token.Type]func(ast.Expression) ast.Expression)
	p.infixParseFns[token.PLUS] = p.parseBinaryExpression
	p.infixParseFns[token.MINUS] = p.parseBinaryExpression
	p.infixParseFns[token.MULTIPLY] = p.parseBinaryExpression
	p.infixParseFns[token.DIVIDE] = p.parseBinaryExpression
	p.infixParseFns[token.MODULO] = p.parseBinaryExpression
	p.infixParseFns[token.EQ] = p.parseBinaryExpression
	p.infixParseFns[token.NOT_EQ] = p.parseBinaryExpression
	p.infixParseFns[token.EQ_STRICT] = p.parseBinaryExpression
	p.infixParseFns[token.NOT_EQ_STRICT] = p.parseBinaryExpression
	p.infixParseFns[token.LT] = p.parseBinaryExpression
	p.infixParseFns[token.GT] = p.parseBinaryExpression
	p.infixParseFns[token.LTE] = p.parseBinaryExpression
	p.infixParseFns[token.GTE] = p.parseBinaryExpression
	p.infixParseFns[token.AND] = p.parseBinaryExpression
	p.infixParseFns[token.OR] = p.parseBinaryExpression
	p.infixParseFns[token.ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[token.LPAREN] = p.parseCallExpression
	p.infixParseFns[token.DOT] = p.parseMemberExpression
	p.infixParseFns[token.LBRACKET] = p.parseComputedMemberExpression

	// Read two tokens, so CurrentToken and PeekToken are both set
	p.NextToken()
	p.NextToken()

	return p
}

// ParseProgram parses the entire program and returns the root AST node
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.CurrentToken.Type != token.EOF {
		stmt := p.parseStatement(p)
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}

	return program
}

// NextToken advances both CurrentToken and PeekToken
func (p *Parser) NextToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

// Errors returns the list of parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// AddError adds an error message to the parser's error list
func (p *Parser) AddError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Line %d, Col %d: %s",
		p.CurrentToken.Line, p.CurrentToken.Column, msg))
}

// ExpectToken checks if PeekToken is of expected type and advances if so
func (p *Parser) ExpectToken(t token.Type) bool {
	if p.PeekToken.Type == t {
		p.NextToken()
		return true
	}
	p.AddError(fmt.Sprintf("expected %s, got %s", t, p.PeekToken.Type))
	return false
}

// peekPrecedence returns the precedence of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.PeekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// currentPrecedence returns the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.CurrentToken.Type]; ok {
		return p
	}
	return LOWEST
}
