// Package parser provides syntax analysis functionality for the XJS language.
// It builds an Abstract Syntax Tree (AST) from tokens provided by the lexer.
package parser

import (
	"fmt"
	"strconv"

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

	currentToken token.Token
	peekToken    token.Token

	errors []string

	prefixParseFns map[token.Type]func() ast.Expression
	infixParseFns  map[token.Type]func(ast.Expression) ast.Expression
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
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

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances both currentToken and peekToken
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Errors returns the list of parsing errors
func (p *Parser) Errors() []string {
	return p.errors
}

// addError adds an error message to the parser's error list
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Line %d, Col %d: %s",
		p.currentToken.Line, p.currentToken.Column, msg))
}

// expectToken checks if peekToken is of expected type and advances if so
func (p *Parser) expectToken(t token.Type) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.addError(fmt.Sprintf("expected %s, got %s", t, p.peekToken.Type))
	return false
}

// peekPrecedence returns the precedence of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// currentPrecedence returns the precedence of the current token
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram parses the entire program and returns the root AST node
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parses any kind of statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.FUNCTION:
		return p.parseFunctionDeclaration()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses variable declarations
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectToken(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if p.peekToken.Type == token.ASSIGN {
		p.nextToken() // consume =
		p.nextToken() // move to value
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseFunctionDeclaration parses function declarations
func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{Token: p.currentToken}

	if !p.expectToken(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectToken(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.expectToken(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionParameters parses function parameter list
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseReturnStatement parses return statements
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	if p.peekToken.Type != token.SEMICOLON && p.peekToken.Type != token.EOF {
		p.nextToken()
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseIfStatement parses if/else statements
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}

	if !p.expectToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	p.nextToken()
	stmt.ThenBranch = p.parseStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()
		p.nextToken()
		stmt.ElseBranch = p.parseStatement()
	}

	return stmt
}

// parseWhileStatement parses while loops
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.currentToken}

	if !p.expectToken(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Body = p.parseStatement()

	return stmt
}

// parseForStatement parses for loops
func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.currentToken}

	if !p.expectToken(token.LPAREN) {
		return nil
	}

	// Parse init
	if p.peekToken.Type != token.SEMICOLON {
		p.nextToken()
		stmt.Init = p.parseStatement()
	} else {
		p.nextToken() // consume semicolon
	}

	// Parse condition
	if p.peekToken.Type != token.SEMICOLON {
		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectToken(token.SEMICOLON) {
		return nil
	}

	// Parse update
	if p.peekToken.Type != token.RPAREN {
		p.nextToken()
		stmt.Update = p.parseExpression(LOWEST)
	}

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Body = p.parseStatement()

	return stmt
}

// parseBlockStatement parses block statements
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseExpressionStatement parses expression statements
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// parseExpression parses expressions using Pratt parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.addError(fmt.Sprintf("no prefix parse function for %s found", p.currentToken.Type))
		return nil
	}

	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Prefix parse functions
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.currentToken}

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.currentToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.currentToken}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(UNARY)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectToken(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.currentToken}
	obj.Properties = make(map[ast.Expression]ast.Expression)

	if p.peekToken.Type == token.RBRACE {
		p.nextToken()
		return obj
	}

	p.nextToken()

	for {
		key := p.parseExpression(LOWEST)

		if !p.expectToken(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		obj.Properties[key] = value

		if p.peekToken.Type != token.COMMA {
			break
		}
		p.nextToken()
		p.nextToken()
	}

	if !p.expectToken(token.RBRACE) {
		return nil
	}

	return obj
}

// Infix parse functions
func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.AssignmentExpression{
		Token: p.currentToken,
		Left:  left,
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currentToken, Function: fn}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.currentToken,
		Object:   left,
		Computed: false,
	}

	p.nextToken()
	exp.Property = p.parseExpression(MEMBER)

	return exp
}

func (p *Parser) parseComputedMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.currentToken,
		Object:   left,
		Computed: true,
	}

	p.nextToken()
	exp.Property = p.parseExpression(LOWEST)

	if !p.expectToken(token.RBRACKET) {
		return nil
	}

	return exp
}

// parseExpressionList parses a comma-separated list of expressions
func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectToken(end) {
		return nil
	}

	return args
}
