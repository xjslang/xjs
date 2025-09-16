package parser

import (
	"fmt"
	"strconv"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}
	if !p.ExpectToken(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	if p.peekToken.Type == token.ASSIGN {
		p.NextToken() // consume =
		p.NextToken() // move to value
		stmt.Value = p.ParseExpression()
	}
	if p.peekToken.Type == token.SEMICOLON {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) ParseFunctionStatement() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{Token: p.currentToken}
	if !p.ExpectToken(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
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
	if p.peekToken.Type == token.RPAREN {
		p.NextToken()
		return identifiers
	}
	p.NextToken()
	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)
	for p.peekToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}
	if !p.ExpectToken(token.RPAREN) {
		return nil
	}
	return identifiers
}

func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	if p.peekToken.Type != token.SEMICOLON && p.peekToken.Type != token.EOF {
		p.NextToken()
		stmt.ReturnValue = p.ParseExpression()
	}
	if p.peekToken.Type == token.SEMICOLON {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) ParseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}
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
	if p.peekToken.Type == token.ELSE {
		p.NextToken()
		p.NextToken()
		stmt.ElseBranch = p.statementParseFn(p)
	}
	return stmt
}

func (p *Parser) ParseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.currentToken}
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
	stmt := &ast.ForStatement{Token: p.currentToken}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	if p.peekToken.Type != token.SEMICOLON {
		p.NextToken()
		stmt.Init = p.statementParseFn(p)
	} else {
		p.NextToken() // consume semicolon
	}
	if p.peekToken.Type != token.SEMICOLON {
		p.NextToken()
		stmt.Condition = p.ParseExpression()
	}
	if !p.ExpectToken(token.SEMICOLON) {
		return nil
	}
	if p.peekToken.Type != token.RPAREN {
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
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.PushContext(BlockContext)
	defer p.PopContext()
	p.NextToken()
	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
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
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.ParseExpression()
	if p.peekToken.Type == token.SEMICOLON {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) ParsePrefixExpression() ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.AddError(fmt.Sprintf("no prefix parse function for %s found", p.currentToken.Type))
		return nil
	}
	return prefix()
}

func (p *Parser) ParseInfixExpression(left ast.Expression) ast.Expression {
	infix := p.infixParseFns[p.peekToken.Type]
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
	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		left = p.ParseInfixExpression(left)
	}
	return left
}

func (p *Parser) ParseRemainingExpression(left ast.Expression) ast.Expression {
	return p.ParseRemainingExpressionWithPrecedence(left, p.currentExpressionPrecedence)
}

func (p *Parser) ParseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}
	_, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal))
		return nil
	}
	return lit
}

func (p *Parser) ParseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.currentToken}
	_, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as float", p.currentToken.Literal))
		return nil
	}
	return lit
}

func (p *Parser) ParseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) ParseMultiStringLiteral() ast.Expression {
	return &ast.MultiStringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) ParseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
}

func (p *Parser) ParseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.currentToken}
}

func (p *Parser) ParseUnaryExpression() ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.NextToken()
	expression.Right = p.expressionParseFn(p, UNARY)
	return expression
}

func (p *Parser) ParsePostfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.PostfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}
	return expression
}

func (p *Parser) ParseGroupedExpression() ast.Expression {
	p.NextToken()
	exp := p.ParseExpression()
	if !p.ExpectToken(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) ParseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.currentToken}
	array.Elements = p.ParseExpressionList(token.RBRACKET)
	return array
}

func (p *Parser) ParseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.currentToken}
	obj.Properties = make(map[ast.Expression]ast.Expression)
	if p.peekToken.Type == token.RBRACE {
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
		obj.Properties[key] = value
		if p.peekToken.Type != token.COMMA {
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
	fe := &ast.FunctionExpression{Token: p.currentToken}
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
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}
	precedence := p.currentPrecedence()
	p.NextToken()
	expression.Right = p.expressionParseFn(p, precedence)
	return expression
}

func (p *Parser) ParseAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.AssignmentExpression{
		Token: p.currentToken,
		Left:  left,
	}
	p.NextToken()
	expression.Value = p.ParseExpression()
	return expression
}

func (p *Parser) ParseCompoundAssignmentExpression(left ast.Expression) ast.Expression {
	expression := &ast.CompoundAssignmentExpression{
		Token: p.currentToken,
		Left:  left,
	}
	switch p.currentToken.Type {
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
	exp := &ast.CallExpression{Token: p.currentToken, Function: fn}
	exp.Arguments = p.ParseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) ParseMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.currentToken,
		Object:   left,
		Computed: false,
	}
	p.NextToken()
	exp.Property = p.expressionParseFn(p, MEMBER)
	return exp
}

func (p *Parser) ParseComputedMemberExpression(left ast.Expression) ast.Expression {
	exp := &ast.MemberExpression{
		Token:    p.currentToken,
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
	if p.peekToken.Type == end {
		p.NextToken()
		return args
	}
	p.NextToken()
	args = append(args, p.ParseExpression())
	for p.peekToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		args = append(args, p.ParseExpression())
	}
	if !p.ExpectToken(end) {
		return nil
	}
	return args
}
