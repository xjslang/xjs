package parser

import (
	"fmt"
	"strconv"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

// ParseCommentBlock parses one or more consecutive comment tokens into a CommentBlock.
func (p *Parser) ParseCommentBlock() *ast.CommentBlock {
	block := &ast.CommentBlock{
		Comments: []token.Token{p.CurrentToken},
	}

	// Collect consecutive comment tokens
	for p.PeekToken.Type == token.COMMENT {
		p.NextToken()
		block.Comments = append(block.Comments, p.CurrentToken)
	}

	return block
}

// ParseBlankLine parses blank line tokens into a BlankLine statement.
// Multiple consecutive BLANK_LINE tokens are condensed into one.
func (p *Parser) ParseBlankLine() *ast.BlankLine {
	bl := &ast.BlankLine{Token: p.CurrentToken}

	// Consume consecutive blank line tokens
	for p.PeekToken.Type == token.BLANK_LINE {
		p.NextToken()
	}

	return bl
}

func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.CurrentToken}
	if !p.ExpectToken(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	if p.PeekToken.Type == token.ASSIGN {
		p.NextToken() // consume =
		p.NextToken() // move to value
		stmt.Value = p.ParseExpression()
	}
	if !p.ExpectSemicolonASI() {
		return nil
	}
	return stmt
}

// ParseLetExpression parses a let expression (used in for loop initializers).
// Unlike ParseLetStatement, it returns an Expression and doesn't expect a semicolon.
func (p *Parser) ParseLetExpression() ast.Expression {
	expr := &ast.LetExpression{Token: p.CurrentToken}
	if !p.ExpectToken(token.IDENT) {
		return nil
	}
	expr.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	if p.PeekToken.Type == token.ASSIGN {
		p.NextToken() // consume =
		p.NextToken() // move to value
		expr.Value = p.ParseExpression()
	}
	return expr
}

func (p *Parser) ParseFunctionStatement() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{Token: p.CurrentToken}
	if !p.ExpectToken(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	stmt.Parameters = p.ParseFunctionParameters()
	if !p.ExpectToken(token.LBRACE) {
		return nil
	}
	p.PushContext(FunctionContext)
	defer p.PopContext()
	stmt.Body = p.ParseBlockStatement()
	return stmt
}

func (p *Parser) ParseFunctionParameters() []*ast.Identifier {
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

func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.CurrentToken}
	if p.PeekToken.Type != token.SEMICOLON && p.PeekToken.Type != token.EOF && p.PeekToken.Type != token.RBRACE && p.PeekToken.Type != token.COMMENT {
		p.NextToken()
		stmt.ReturnValue = p.ParseExpression()
	}
	if !p.ExpectSemicolonASI() {
		return nil
	}
	return stmt
}

func (p *Parser) ParseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.CurrentToken}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	p.NextToken()
	stmt.Condition = p.ParseExpression()
	if !p.ExpectToken(token.RPAREN) {
		return nil
	}
	p.NextToken()
	stmt.ThenBranch = p.statementParseFn(p)
	if p.PeekToken.Type == token.ELSE {
		p.NextToken()
		p.NextToken()
		stmt.ElseBranch = p.statementParseFn(p)
	}
	return stmt
}

func (p *Parser) ParseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.CurrentToken}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	p.NextToken()
	stmt.Condition = p.ParseExpression()
	if !p.ExpectToken(token.RPAREN) {
		return nil
	}
	p.NextToken()
	stmt.Body = p.statementParseFn(p)
	return stmt
}

func (p *Parser) ParseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.CurrentToken}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	if p.PeekToken.Type != token.SEMICOLON {
		p.NextToken()
		// Parse init as expression: either LetExpression or regular expression
		if p.CurrentToken.Type == token.LET {
			stmt.Init = p.ParseLetExpression()
		} else {
			stmt.Init = p.ParseExpression()
		}
	}
	if !p.ExpectToken(token.SEMICOLON) {
		return nil
	}
	if p.PeekToken.Type != token.SEMICOLON {
		p.NextToken()
		stmt.Condition = p.ParseExpression()
	}
	if !p.ExpectToken(token.SEMICOLON) {
		return nil
	}
	if p.PeekToken.Type != token.RPAREN {
		p.NextToken()
		stmt.Update = p.ParseExpression()
	}
	if !p.ExpectToken(token.RPAREN) {
		return nil
	}
	p.NextToken()
	stmt.Body = p.statementParseFn(p)
	return stmt
}

func (p *Parser) ParseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.CurrentToken}
	block.Statements = []ast.Statement{}
	p.PushContext(BlockContext)
	defer p.PopContext()
	p.NextToken()
	for p.CurrentToken.Type != token.RBRACE && p.CurrentToken.Type != token.EOF {
		stmt := p.statementParseFn(p)
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.NextToken()
	}
	return block
}

func (p *Parser) ParseStatement() ast.Statement {
	return p.statementParseFn(p)
}

func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.CurrentToken}
	stmt.Expression = p.ParseExpression()
	if !p.ExpectSemicolonASI() {
		return nil
	}
	return stmt
}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	prefix := p.prefixParseFns[p.CurrentToken.Type]
	if prefix == nil {
		p.AddError(fmt.Sprintf("unexpected %s", p.CurrentToken.Literal))
		return nil
	}
	return prefix()
}

func (p *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	infix := p.infixParseFns[p.PeekToken.Type]
	if infix == nil {
		return left
	}
	p.NextToken()
	return infix(left)
}

func (p *Parser) ParseExpression() ast.Expression {
	return p.expressionParseFn(p, LOWEST)
}

func (p *Parser) ParseExpressionWithPrecedence(precedence int) ast.Expression {
	return p.expressionParseFn(p, precedence)
}

func (p *Parser) ParseRemainingExpressionWithPrecedence(left ast.Expression, precedence int) ast.Expression {
	for p.PeekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		// Smart semicolon insertion: prevent LPAREN and LBRACKET after newline from continuing expression
		// https://eslint.org/docs/latest/rules/no-unexpected-multiline
		if p.smartSemicolons && p.PeekToken.AfterNewline {
			switch p.PeekToken.Type {
			case token.LPAREN, token.LBRACKET:
				// These tokens after a newline should not continue the expression
				return left
			}
		}
		// Stop if we hit a comment or blank line
		if p.PeekToken.Type == token.COMMENT || p.PeekToken.Type == token.BLANK_LINE {
			return left
		}
		left = p.ParseInfixExpression(left)
	}
	return left
}

func (p *Parser) ParseRemainingExpression(left ast.Expression) ast.Expression {
	return p.ParseRemainingExpressionWithPrecedence(left, p.currentExpressionPrecedence)
}

func (p *Parser) ParseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.CurrentToken}
	_, err := strconv.ParseInt(p.CurrentToken.Literal, 0, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as integer", p.CurrentToken.Literal))
		return nil
	}
	return lit
}

func (p *Parser) ParseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.CurrentToken}
	_, err := strconv.ParseFloat(p.CurrentToken.Literal, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as float", p.CurrentToken.Literal))
		return nil
	}
	return lit
}

func (p *Parser) ParseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func (p *Parser) ParseMultiStringLiteral() ast.Expression {
	return &ast.MultiStringLiteral{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func (p *Parser) ParseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.CurrentToken, Value: p.CurrentToken.Type == token.TRUE}
}

func (p *Parser) ParseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.CurrentToken}
}

func (p *Parser) ParseUnaryExpression() ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.CurrentToken,
		Operator: p.CurrentToken.Literal,
	}
	p.NextToken()
	expression.Right = p.expressionParseFn(p, UNARY)
	return expression
}

func (p *Parser) ParsePostfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.PostfixExpression{
		Token:    p.CurrentToken,
		Left:     left,
		Operator: p.CurrentToken.Literal,
	}
	return expression
}

func (p *Parser) ParseGroupedExpression() ast.Expression {
	tok := p.CurrentToken
	p.NextToken()
	exp := p.ParseExpression()
	if !p.ExpectToken(token.RPAREN) {
		return nil
	}
	return &ast.GroupedExpression{
		Token:      tok,
		Expression: exp,
	}
}

func (p *Parser) ParseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.CurrentToken}
	array.Elements = p.ParseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) ParseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.CurrentToken}
	obj.Properties = []ast.ObjectProperty{}
	if p.PeekToken.Type == token.RBRACE {
		p.NextToken()
		return obj
	}
	p.NextToken()
	for {
		key := p.ParseExpression()
		if !p.ExpectToken(token.COLON) {
			return nil
		}
		p.NextToken()
		value := p.ParseExpression()
		obj.Properties = append(obj.Properties, ast.ObjectProperty{
			Key:   key,
			Value: value,
		})
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

func (p *Parser) ParseFunctionExpression() ast.Expression {
	fe := &ast.FunctionExpression{Token: p.CurrentToken}
	// Check for optional function name
	if p.PeekToken.Type == token.IDENT {
		p.NextToken()
		fe.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	fe.Parameters = p.ParseFunctionParameters()
	if !p.ExpectToken(token.LBRACE) {
		return nil
	}
	p.PushContext(FunctionContext)
	defer p.PopContext()
	fe.Body = p.ParseBlockStatement()
	return fe
}

func (p *Parser) ParseBinaryExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.CurrentToken,
		Left:     left,
		Operator: p.CurrentToken.Literal,
	}
	precedence := p.currentPrecedence()
	p.NextToken()
	expression.Right = p.expressionParseFn(p, precedence)
	return expression
}

func (p *Parser) ParseAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.AssignmentExpression{
		Token: p.CurrentToken,
		Left:  left,
	}
	p.NextToken()
	expression.Value = p.ParseExpression()
	return expression
}

func (p *Parser) ParseCompoundAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.CompoundAssignmentExpression{
		Token: p.CurrentToken,
		Left:  left,
	}
	switch p.CurrentToken.Type {
	case token.PLUS_ASSIGN:
		expression.Operator = "+"
	case token.MINUS_ASSIGN:
		expression.Operator = "-"
	}
	p.NextToken()
	expression.Value = p.ParseExpression()
	return expression
}

func (p *Parser) ParseCallExpression(fn ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.CurrentToken, Function: fn}
	exp.Arguments = p.ParseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) ParseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.CurrentToken,
		Object:   left,
		Computed: false,
	}
	p.NextToken()
	exp.Property = p.expressionParseFn(p, MEMBER)
	return exp
}

func (p *Parser) ParseComputedMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.CurrentToken,
		Object:   left,
		Computed: true,
	}
	p.NextToken()
	exp.Property = p.ParseExpression()
	if !p.ExpectToken(token.RBRACKET) {
		return nil
	}
	return exp
}

func (p *Parser) ParseExpressionList(end token.Type) []ast.Expression {
	args := []ast.Expression{}
	if p.PeekToken.Type == end {
		p.NextToken()
		return args
	}
	p.NextToken()
	args = append(args, p.ParseExpression())
	for p.PeekToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		args = append(args, p.ParseExpression())
	}
	if !p.ExpectToken(end) {
		return nil
	}
	return args
}

// shouldInsertSemicolon determines if ASI (Automatic Semicolon Insertion) should occur
// based on JavaScript's ASI rules:
// 1. There is a line terminator between two tokens that cannot be part of the same statement
// 2. The next token is '}' (closing brace)
// 3. The next token is EOF
func (p *Parser) shouldInsertSemicolon() bool {
	if p.PeekToken.Type == token.EOF {
		return true
	}
	if p.PeekToken.Type == token.RBRACE {
		return true
	}
	// Comments and blank lines after a statement should trigger ASI
	if p.PeekToken.Type == token.COMMENT || p.PeekToken.Type == token.BLANK_LINE {
		return true
	}
	if !p.PeekToken.AfterNewline {
		return false
	}
	switch p.PeekToken.Type {
	case token.DOT: // obj.prop
	case token.LBRACE: // object literal (e.g., return { ... })
	case token.PLUS: // + (binary)
	case token.MINUS: // - (binary)
	case token.MULTIPLY: // *
	case token.DIVIDE: // /
	case token.MODULO: // %
	case token.LT: // <
	case token.GT: // >
	case token.LTE: // <=
	case token.GTE: // >=
	case token.EQ: // ==
	case token.NOT_EQ: // !=
	case token.AND: // &&
	case token.OR: // ||
	case token.ASSIGN: // =
	case token.PLUS_ASSIGN: // +=
	case token.MINUS_ASSIGN: // -=
		return false
	}
	return true
}
