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

	CurrentToken token.Token
	PeekToken    token.Token

	parseStatement    func(p *Parser) ast.Statement
	parseLetStatement func(p *Parser) *ast.LetStatement

	prefixParseFns map[token.Type]func() ast.Expression
	infixParseFns  map[token.Type]func(ast.Expression) ast.Expression

	errors []string
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:             l,
		errors:            []string{},
		parseStatement:    baseParseStatement,
		parseLetStatement: baseParseLetStatement,
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

// parseFunctionDeclaration parses function declarations
func (p *Parser) parseFunctionDeclaration() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{Token: p.CurrentToken}

	if !p.ExpectToken(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}

	if !p.ExpectToken(token.LPAREN) {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.ExpectToken(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseFunctionParameters parses function parameter list
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.PeekToken.Type == token.RPAREN {
		p.NextToken()
		return identifiers
	}

	p.NextToken()

	ident := &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.PeekToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		ident := &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.ExpectToken(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseReturnStatement parses return statements
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.CurrentToken}

	if p.PeekToken.Type != token.SEMICOLON && p.PeekToken.Type != token.EOF {
		p.NextToken()
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.PeekToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

// parseIfStatement parses if/else statements
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.CurrentToken}

	if !p.ExpectToken(token.LPAREN) {
		return nil
	}

	p.NextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.ExpectToken(token.RPAREN) {
		return nil
	}

	p.NextToken()
	stmt.ThenBranch = p.parseStatement(p)

	if p.PeekToken.Type == token.ELSE {
		p.NextToken()
		p.NextToken()
		stmt.ElseBranch = p.parseStatement(p)
	}

	return stmt
}

// parseWhileStatement parses while loops
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.CurrentToken}

	if !p.ExpectToken(token.LPAREN) {
		return nil
	}

	p.NextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.ExpectToken(token.RPAREN) {
		return nil
	}

	p.NextToken()
	stmt.Body = p.parseStatement(p)

	return stmt
}

// parseForStatement parses for loops
func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.CurrentToken}

	if !p.ExpectToken(token.LPAREN) {
		return nil
	}

	// Parse init
	if p.PeekToken.Type != token.SEMICOLON {
		p.NextToken()
		stmt.Init = p.parseStatement(p)
	} else {
		p.NextToken() // consume semicolon
	}

	// Parse condition
	if p.PeekToken.Type != token.SEMICOLON {
		p.NextToken()
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if !p.ExpectToken(token.SEMICOLON) {
		return nil
	}

	// Parse update
	if p.PeekToken.Type != token.RPAREN {
		p.NextToken()
		stmt.Update = p.parseExpression(LOWEST)
	}

	if !p.ExpectToken(token.RPAREN) {
		return nil
	}

	p.NextToken()
	stmt.Body = p.parseStatement(p)

	return stmt
}

// parseBlockStatement parses block statements
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.CurrentToken}
	block.Statements = []ast.Statement{}

	p.NextToken()

	for p.CurrentToken.Type != token.RBRACE && p.CurrentToken.Type != token.EOF {
		stmt := p.parseStatement(p)
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}

	return block
}

// parseExpressionStatement parses expression statements
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.CurrentToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.PeekToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}

// parseExpression parses expressions using Pratt parsing
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.CurrentToken.Type]
	if prefix == nil {
		p.AddError(fmt.Sprintf("no prefix parse function for %s found", p.CurrentToken.Type))
		return nil
	}

	leftExp := prefix()

	for p.PeekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.PeekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.NextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Prefix parse functions
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.CurrentToken}

	value, err := strconv.ParseInt(p.CurrentToken.Literal, 0, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as integer", p.CurrentToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.CurrentToken}

	value, err := strconv.ParseFloat(p.CurrentToken.Literal, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as float", p.CurrentToken.Literal))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.CurrentToken, Value: p.CurrentToken.Type == token.TRUE}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.CurrentToken}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.CurrentToken,
		Operator: p.CurrentToken.Literal,
	}

	p.NextToken()
	expression.Right = p.parseExpression(UNARY)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken()

	exp := p.parseExpression(LOWEST)

	if !p.ExpectToken(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.CurrentToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.CurrentToken}
	obj.Properties = make(map[ast.Expression]ast.Expression)

	if p.PeekToken.Type == token.RBRACE {
		p.NextToken()
		return obj
	}

	p.NextToken()

	for {
		key := p.parseExpression(LOWEST)

		if !p.ExpectToken(token.COLON) {
			return nil
		}

		p.NextToken()
		value := p.parseExpression(LOWEST)

		obj.Properties[key] = value

		if p.PeekToken.Type != token.COMMA {
			break
		}
		p.NextToken()
		p.NextToken()
	}

	if !p.ExpectToken(token.RBRACE) {
		return nil
	}

	return obj
}

// Infix parse functions
func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.CurrentToken,
		Left:     left,
		Operator: p.CurrentToken.Literal,
	}

	precedence := p.currentPrecedence()
	p.NextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.AssignmentExpression{
		Token: p.CurrentToken,
		Left:  left,
	}

	p.NextToken()
	expression.Value = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.CurrentToken, Function: fn}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.CurrentToken,
		Object:   left,
		Computed: false,
	}

	p.NextToken()
	exp.Property = p.parseExpression(MEMBER)

	return exp
}

func (p *Parser) parseComputedMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.CurrentToken,
		Object:   left,
		Computed: true,
	}

	p.NextToken()
	exp.Property = p.parseExpression(LOWEST)

	if !p.ExpectToken(token.RBRACKET) {
		return nil
	}

	return exp
}

// parseExpressionList parses a comma-separated list of expressions
func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	args := []ast.Expression{}

	if p.PeekToken.Type == end {
		p.NextToken()
		return args
	}

	p.NextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.PeekToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.ExpectToken(end) {
		return nil
	}

	return args
}
