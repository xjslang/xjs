package main

import (
	"fmt"
	"strconv"
)

// AST Node interfaces
type Node interface {
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program represents the root of every AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out string
	for _, stmt := range p.Statements {
		out += stmt.String()
	}
	return out
}

// Statements
type LetStatement struct {
	Token Token // the LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) String() string {
	out := "let " + ls.Name.String()
	if ls.Value != nil {
		out += " = " + ls.Value.String()
	}
	out += ";"
	return out
}

type ReturnStatement struct {
	Token       Token // the RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) String() string {
	out := "return"
	if rs.ReturnValue != nil {
		out += " " + rs.ReturnValue.String()
	}
	out += ";"
	return out
}

type ExpressionStatement struct {
	Token      Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type FunctionDeclaration struct {
	Token      Token // the FUNCTION token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fd *FunctionDeclaration) statementNode() {}
func (fd *FunctionDeclaration) String() string {
	out := "function " + fd.Name.String() + "("
	for i, param := range fd.Parameters {
		if i > 0 {
			out += ", "
		}
		out += param.String()
	}
	out += ") " + fd.Body.String()
	return out
}

type BlockStatement struct {
	Token      Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) String() string {
	out := "{"
	for _, stmt := range bs.Statements {
		out += stmt.String()
	}
	out += "}"
	return out
}

type IfStatement struct {
	Token     Token // the IF token
	Condition Expression
	ThenBranch Statement
	ElseBranch Statement // can be nil
}

func (ifs *IfStatement) statementNode() {}
func (ifs *IfStatement) String() string {
	out := "if (" + ifs.Condition.String() + ") " + ifs.ThenBranch.String()
	if ifs.ElseBranch != nil {
		out += " else " + ifs.ElseBranch.String()
	}
	return out
}

type WhileStatement struct {
	Token     Token // the WHILE token
	Condition Expression
	Body      Statement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) String() string {
	return "while (" + ws.Condition.String() + ") " + ws.Body.String()
}

type ForStatement struct {
	Token     Token // the FOR token
	Init      Statement   // can be nil
	Condition Expression  // can be nil
	Update    Expression  // can be nil
	Body      Statement
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	out := "for ("
	if fs.Init != nil {
		out += fs.Init.String()
	}
	out += "; "
	if fs.Condition != nil {
		out += fs.Condition.String()
	}
	out += "; "
	if fs.Update != nil {
		out += fs.Update.String()
	}
	out += ") " + fs.Body.String()
	return out
}

// Expressions
type Identifier struct {
	Token Token // the IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string { return i.Value }

type IntegerLiteral struct {
	Token Token // the INT token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) String() string { return il.Token.Literal }

type FloatLiteral struct {
	Token Token // the FLOAT token
	Value float64
}

func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) String() string { return fl.Token.Literal }

type StringLiteral struct {
	Token Token // the STRING token
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) String() string { return "\"" + sl.Value + "\"" }

type BooleanLiteral struct {
	Token Token // the TRUE or FALSE token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) String() string { return bl.Token.Literal }

type NullLiteral struct {
	Token Token // the NULL token
}

func (nl *NullLiteral) expressionNode() {}
func (nl *NullLiteral) String() string { return "null" }

type BinaryExpression struct {
	Token    Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}
func (be *BinaryExpression) String() string {
	return "(" + be.Left.String() + " " + be.Operator + " " + be.Right.String() + ")"
}

type UnaryExpression struct {
	Token    Token // the operator token
	Operator string
	Right    Expression
}

func (ue *UnaryExpression) expressionNode() {}
func (ue *UnaryExpression) String() string {
	return "(" + ue.Operator + ue.Right.String() + ")"
}

type CallExpression struct {
	Token     Token // the ( token
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) String() string {
	out := ce.Function.String() + "("
	for i, arg := range ce.Arguments {
		if i > 0 {
			out += ", "
		}
		out += arg.String()
	}
	out += ")"
	return out
}

type MemberExpression struct {
	Token    Token // the . or [ token
	Object   Expression
	Property Expression
	Computed bool // true for obj[prop], false for obj.prop
}

func (me *MemberExpression) expressionNode() {}
func (me *MemberExpression) String() string {
	if me.Computed {
		return me.Object.String() + "[" + me.Property.String() + "]"
	}
	return me.Object.String() + "." + me.Property.String()
}

type AssignmentExpression struct {
	Token Token // the = token
	Left  Expression
	Value Expression
}

func (ae *AssignmentExpression) expressionNode() {}
func (ae *AssignmentExpression) String() string {
	return ae.Left.String() + " = " + ae.Value.String()
}

type ArrayLiteral struct {
	Token    Token // the [ token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) String() string {
	out := "["
	for i, elem := range al.Elements {
		if i > 0 {
			out += ", "
		}
		out += elem.String()
	}
	out += "]"
	return out
}

type ObjectLiteral struct {
	Token      Token // the { token
	Properties map[Expression]Expression
}

func (ol *ObjectLiteral) expressionNode() {}
func (ol *ObjectLiteral) String() string {
	out := "{"
	i := 0
	for key, value := range ol.Properties {
		if i > 0 {
			out += ", "
		}
		out += key.String() + ": " + value.String()
		i++
	}
	out += "}"
	return out
}

// Operator precedence
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

// Token precedences
var precedences = map[TokenType]int{
	ASSIGN:    ASSIGNMENT,
	OR:        LOGICAL_OR,
	AND:       LOGICAL_AND,
	EQ:        EQUALITY,
	NOT_EQ:    EQUALITY,
	EQ_STRICT: EQUALITY,
	NOT_EQ_STRICT: EQUALITY,
	LT:        COMPARISON,
	GT:        COMPARISON,
	LTE:       COMPARISON,
	GTE:       COMPARISON,
	PLUS:      SUM,
	MINUS:     SUM,
	MULTIPLY:  PRODUCT,
	DIVIDE:    PRODUCT,
	MODULO:    PRODUCT,
	LPAREN:    CALL,
	DOT:       MEMBER,
	LBRACKET:  MEMBER,
}

// Parser structure
type Parser struct {
	lexer *Lexer
	
	currentToken Token
	peekToken    Token
	
	errors []string
	
	prefixParseFns map[TokenType]func() Expression
	infixParseFns  map[TokenType]func(Expression) Expression
}

// NewParser creates a new parser instance
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}
	
	// Initialize prefix parse functions
	p.prefixParseFns = make(map[TokenType]func() Expression)
	p.prefixParseFns[IDENT] = p.parseIdentifier
	p.prefixParseFns[INT] = p.parseIntegerLiteral
	p.prefixParseFns[FLOAT] = p.parseFloatLiteral
	p.prefixParseFns[STRING] = p.parseStringLiteral
	p.prefixParseFns[TRUE] = p.parseBooleanLiteral
	p.prefixParseFns[FALSE] = p.parseBooleanLiteral
	p.prefixParseFns[NULL] = p.parseNullLiteral
	p.prefixParseFns[NOT] = p.parseUnaryExpression
	p.prefixParseFns[MINUS] = p.parseUnaryExpression
	p.prefixParseFns[INCREMENT] = p.parseUnaryExpression
	p.prefixParseFns[DECREMENT] = p.parseUnaryExpression
	p.prefixParseFns[LPAREN] = p.parseGroupedExpression
	p.prefixParseFns[LBRACKET] = p.parseArrayLiteral
	p.prefixParseFns[LBRACE] = p.parseObjectLiteral
	
	// Initialize infix parse functions
	p.infixParseFns = make(map[TokenType]func(Expression) Expression)
	p.infixParseFns[PLUS] = p.parseBinaryExpression
	p.infixParseFns[MINUS] = p.parseBinaryExpression
	p.infixParseFns[MULTIPLY] = p.parseBinaryExpression
	p.infixParseFns[DIVIDE] = p.parseBinaryExpression
	p.infixParseFns[MODULO] = p.parseBinaryExpression
	p.infixParseFns[EQ] = p.parseBinaryExpression
	p.infixParseFns[NOT_EQ] = p.parseBinaryExpression
	p.infixParseFns[EQ_STRICT] = p.parseBinaryExpression
	p.infixParseFns[NOT_EQ_STRICT] = p.parseBinaryExpression
	p.infixParseFns[LT] = p.parseBinaryExpression
	p.infixParseFns[GT] = p.parseBinaryExpression
	p.infixParseFns[LTE] = p.parseBinaryExpression
	p.infixParseFns[GTE] = p.parseBinaryExpression
	p.infixParseFns[AND] = p.parseBinaryExpression
	p.infixParseFns[OR] = p.parseBinaryExpression
	p.infixParseFns[ASSIGN] = p.parseAssignmentExpression
	p.infixParseFns[LPAREN] = p.parseCallExpression
	p.infixParseFns[DOT] = p.parseMemberExpression
	p.infixParseFns[LBRACKET] = p.parseComputedMemberExpression
	
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

// Errors returns parser errors
func (p *Parser) Errors() []string {
	return p.errors
}

// addError adds an error message
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("Line %d, Col %d: %s", 
		p.currentToken.Line, p.currentToken.Column, msg))
}

// expectToken checks if peekToken is of expected type and advances if so
func (p *Parser) expectToken(t TokenType) bool {
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

// ParseProgram parses the entire program
func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}
	
	for p.currentToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	
	return program
}

// parseStatement parses a statement
func (p *Parser) parseStatement() Statement {
	switch p.currentToken.Type {
	case LET:
		return p.parseLetStatement()
	case FUNCTION:
		return p.parseFunctionDeclaration()
	case RETURN:
		return p.parseReturnStatement()
	case IF:
		return p.parseIfStatement()
	case WHILE:
		return p.parseWhileStatement()
	case FOR:
		return p.parseForStatement()
	case LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement parses let statements
func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.currentToken}
	
	if !p.expectToken(IDENT) {
		return nil
	}
	
	stmt.Name = &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	
	if p.peekToken.Type == ASSIGN {
		p.nextToken() // consume =
		p.nextToken() // move to value
		stmt.Value = p.parseExpression(LOWEST)
	}
	
	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}
	
	return stmt
}

// parseFunctionDeclaration parses function declarations
func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
	stmt := &FunctionDeclaration{Token: p.currentToken}
	
	if !p.expectToken(IDENT) {
		return nil
	}
	
	stmt.Name = &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	
	if !p.expectToken(LPAREN) {
		return nil
	}
	
	stmt.Parameters = p.parseFunctionParameters()
	
	if !p.expectToken(LBRACE) {
		return nil
	}
	
	stmt.Body = p.parseBlockStatement()
	
	return stmt
}

// parseFunctionParameters parses function parameters
func (p *Parser) parseFunctionParameters() []*Identifier {
	identifiers := []*Identifier{}
	
	if p.peekToken.Type == RPAREN {
		p.nextToken()
		return identifiers
	}
	
	p.nextToken()
	
	ident := &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)
	
	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}
	
	if !p.expectToken(RPAREN) {
		return nil
	}
	
	return identifiers
}

// parseReturnStatement parses return statements
func (p *Parser) parseReturnStatement() *ReturnStatement {
	stmt := &ReturnStatement{Token: p.currentToken}
	
	if p.peekToken.Type != SEMICOLON && p.peekToken.Type != EOF {
		p.nextToken()
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}
	
	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}
	
	return stmt
}

// parseIfStatement parses if statements
func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Token: p.currentToken}
	
	if !p.expectToken(LPAREN) {
		return nil
	}
	
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	
	if !p.expectToken(RPAREN) {
		return nil
	}
	
	p.nextToken()
	stmt.ThenBranch = p.parseStatement()
	
	if p.peekToken.Type == ELSE {
		p.nextToken()
		p.nextToken()
		stmt.ElseBranch = p.parseStatement()
	}
	
	return stmt
}

// parseWhileStatement parses while statements
func (p *Parser) parseWhileStatement() *WhileStatement {
	stmt := &WhileStatement{Token: p.currentToken}
	
	if !p.expectToken(LPAREN) {
		return nil
	}
	
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	
	if !p.expectToken(RPAREN) {
		return nil
	}
	
	p.nextToken()
	stmt.Body = p.parseStatement()
	
	return stmt
}

// parseForStatement parses for statements
func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{Token: p.currentToken}
	
	if !p.expectToken(LPAREN) {
		return nil
	}
	
	// Parse init
	if p.peekToken.Type != SEMICOLON {
		p.nextToken()
		stmt.Init = p.parseStatement()
	} else {
		p.nextToken() // consume semicolon
	}
	
	// Parse condition
	if p.peekToken.Type != SEMICOLON {
		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)
	}
	
	if !p.expectToken(SEMICOLON) {
		return nil
	}
	
	// Parse update
	if p.peekToken.Type != RPAREN {
		p.nextToken()
		stmt.Update = p.parseExpression(LOWEST)
	}
	
	if !p.expectToken(RPAREN) {
		return nil
	}
	
	p.nextToken()
	stmt.Body = p.parseStatement()
	
	return stmt
}

// parseBlockStatement parses block statements
func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.currentToken}
	block.Statements = []Statement{}
	
	p.nextToken()
	
	for p.currentToken.Type != RBRACE && p.currentToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	
	return block
}

// parseExpressionStatement parses expression statements
func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)
	
	if p.peekToken.Type == SEMICOLON {
		p.nextToken()
	}
	
	return stmt
}

// parseExpression parses expressions using Pratt parsing
func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.addError(fmt.Sprintf("no prefix parse function for %s found", p.currentToken.Type))
		return nil
	}
	
	leftExp := prefix()
	
	for p.peekToken.Type != SEMICOLON && precedence < p.peekPrecedence() {
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
func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.currentToken}
	
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal))
		return nil
	}
	
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expression {
	lit := &FloatLiteral{Token: p.currentToken}
	
	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.currentToken.Literal))
		return nil
	}
	
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == TRUE}
}

func (p *Parser) parseNullLiteral() Expression {
	return &NullLiteral{Token: p.currentToken}
}

func (p *Parser) parseUnaryExpression() Expression {
	expression := &UnaryExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	
	p.nextToken()
	expression.Right = p.parseExpression(UNARY)
	
	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()
	
	exp := p.parseExpression(LOWEST)
	
	if !p.expectToken(RPAREN) {
		return nil
	}
	
	return exp
}

func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{Token: p.currentToken}
	array.Elements = p.parseExpressionList(RBRACKET)
	return array
}

func (p *Parser) parseObjectLiteral() Expression {
	obj := &ObjectLiteral{Token: p.currentToken}
	obj.Properties = make(map[Expression]Expression)
	
	if p.peekToken.Type == RBRACE {
		p.nextToken()
		return obj
	}
	
	p.nextToken()
	
	for {
		key := p.parseExpression(LOWEST)
		
		if !p.expectToken(COLON) {
			return nil
		}
		
		p.nextToken()
		value := p.parseExpression(LOWEST)
		
		obj.Properties[key] = value
		
		if p.peekToken.Type != COMMA {
			break
		}
		p.nextToken()
		p.nextToken()
	}
	
	if !p.expectToken(RBRACE) {
		return nil
	}
	
	return obj
}

// Infix parse functions
func (p *Parser) parseBinaryExpression(left Expression) Expression {
	expression := &BinaryExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}
	
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	
	return expression
}

func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	expression := &AssignmentExpression{
		Token: p.currentToken,
		Left:  left,
	}
	
	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)
	
	return expression
}

func (p *Parser) parseCallExpression(fn Expression) Expression {
	exp := &CallExpression{Token: p.currentToken, Function: fn}
	exp.Arguments = p.parseExpressionList(RPAREN)
	return exp
}

func (p *Parser) parseMemberExpression(left Expression) Expression {
	exp := &MemberExpression{
		Token:    p.currentToken,
		Object:   left,
		Computed: false,
	}
	
	p.nextToken()
	exp.Property = p.parseExpression(MEMBER)
	
	return exp
}

func (p *Parser) parseComputedMemberExpression(left Expression) Expression {
	exp := &MemberExpression{
		Token:    p.currentToken,
		Object:   left,
		Computed: true,
	}
	
	p.nextToken()
	exp.Property = p.parseExpression(LOWEST)
	
	if !p.expectToken(RBRACKET) {
		return nil
	}
	
	return exp
}

// parseExpressionList parses a list of expressions separated by commas
func (p *Parser) parseExpressionList(end TokenType) []Expression {
	args := []Expression{}
	
	if p.peekToken.Type == end {
		p.nextToken()
		return args
	}
	
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	
	for p.peekToken.Type == COMMA {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	
	if !p.expectToken(end) {
		return nil
	}
	
	return args
}