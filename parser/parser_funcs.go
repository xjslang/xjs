package parser

import (
	"fmt"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) ParseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.CurrentToken}
	if !p.ExpectToken(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
	if p.PeekToken.Type == token.ASSIGN {
		p.NextToken() // consume =
		p.NextToken() // move to value
		stmt.Value = p.ParseExpression(LOWEST)
	}
	if p.PeekToken.Type == token.SEMICOLON {
		p.NextToken()
	}
	return stmt
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
	if p.PeekToken.Type != token.SEMICOLON && p.PeekToken.Type != token.EOF {
		p.NextToken()
		stmt.ReturnValue = p.ParseExpression(LOWEST)
	}
	if p.PeekToken.Type == token.SEMICOLON {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) ParseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.CurrentToken}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	p.NextToken()
	stmt.Condition = p.ParseExpression(LOWEST)
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
	stmt.Condition = p.ParseExpression(LOWEST)
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
		stmt.Init = p.statementParseFn(p)
	} else {
		p.NextToken() // consume semicolon
	}
	if p.PeekToken.Type != token.SEMICOLON {
		p.NextToken()
		stmt.Condition = p.ParseExpression(LOWEST)
	}
	if !p.ExpectToken(token.SEMICOLON) {
		return nil
	}
	if p.PeekToken.Type != token.RPAREN {
		p.NextToken()
		stmt.Update = p.ParseExpression(LOWEST)
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
	return p.parseExpressionStatement(p)
}

func (p *Parser) ParseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.CurrentToken.Type]
	if prefix == nil {
		p.AddError(fmt.Sprintf("no prefix parse function for %s found", p.CurrentToken.Type))
		return nil
	}
	leftExp := prefix(p)
	for p.PeekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.PeekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.NextToken()
		leftExp = infix(p, leftExp)
	}
	return leftExp
}

func (p *Parser) ParseIdentifier() ast.Expression {
	return baseParseIdentifier(p)
}

func (p *Parser) ParseIntegerLiteral() ast.Expression {
	return baseParseIntegerLiteral(p)
}

func (p *Parser) ParseFloatLiteral() ast.Expression {
	return baseParseFloatLiteral(p)
}

func (p *Parser) ParseStringLiteral() ast.Expression {
	return baseParseStringLiteral(p)
}

func (p *Parser) ParseBooleanLiteral() ast.Expression {
	return baseParseBooleanLiteral(p)
}

func (p *Parser) ParseNullLiteral() ast.Expression {
	return baseParseNullLiteral(p)
}

func (p *Parser) ParseUnaryExpression() ast.Expression {
	return baseParseUnaryExpression(p)
}

func (p *Parser) ParseGroupedExpression() ast.Expression {
	return baseParseGroupedExpression(p)
}

func (p *Parser) ParseArrayLiteral() ast.Expression {
	return baseParseArrayLiteral(p)
}

func (p *Parser) ParseObjectLiteral() ast.Expression {
	return baseParseObjectLiteral(p)
}

func (p *Parser) ParseFunctionExpression() ast.Expression {
	return baseParseFunctionExpression(p)
}

func (p *Parser) ParseBinaryExpression(left ast.Expression) ast.Expression {
	return baseParseBinaryExpression(p, left)
}

func (p *Parser) ParseAssignmentExpression(left ast.Expression) ast.Expression {
	return baseParseAssignmentExpression(p, left)
}

func (p *Parser) ParseCallExpression(fn ast.Expression) ast.Expression {
	return baseParseCallExpression(p, fn)
}

func (p *Parser) ParseMemberExpression(left ast.Expression) ast.Expression {
	return baseParseMemberExpression(p, left)
}

func (p *Parser) ParseComputedMemberExpression(left ast.Expression) ast.Expression {
	return baseParseComputedMemberExpression(p, left)
}

func (p *Parser) ParseExpressionList(end token.Type) []ast.Expression {
	args := []ast.Expression{}
	if p.PeekToken.Type == end {
		p.NextToken()
		return args
	}
	p.NextToken()
	args = append(args, p.ParseExpression(LOWEST))
	for p.PeekToken.Type == token.COMMA {
		p.NextToken()
		p.NextToken()
		args = append(args, p.ParseExpression(LOWEST))
	}
	if !p.ExpectToken(end) {
		return nil
	}
	return args
}
