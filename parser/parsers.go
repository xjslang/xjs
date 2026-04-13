package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) parseBody() *ast.BlockStatement {
	bodyStmt := &ast.BlockStatement{}
	for stmt, err := p.parseStatement(); err == nil; stmt, err = p.parseStatement() {
		bodyStmt.Statements = append(bodyStmt.Statements, stmt)
	}
	return bodyStmt
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.CurrentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.FUNCTION:
		return p.parseFunction()
	}
	return nil, errors.New("Unexpected token " + p.CurrentToken.Literal)
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{}
	p.AdvanceToken() // consume token.LET
	ident, err := p.Expect(token.IDENT)
	if err != nil {
		return stmt, err
	}
	stmt.Name = &ast.Identifier{Value: ident.Literal}
	if _, err := p.Expect(token.ASSIGN); err != nil {
		return stmt, err
	}
	stmt.Value = p.ParseExpression()
	if _, err := p.Expect(token.SEMICOLON); err != nil {
		return stmt, err
	}
	return stmt, nil
}

func (p *Parser) ParseExpression() ast.Expression {
	switch p.CurrentToken.Type {
	case token.NUMBER:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.IntegerLiteral{Value: val}
	case token.STRING:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.StringLiteral{Value: val}
	}
	return nil
}

func (p *Parser) parseFunction() (*ast.FunctionStatement, error) {
	stmt := &ast.FunctionStatement{}
	p.AdvanceToken() // consume token.FUNCTION
	if _, err := p.Expect(token.IDENT); err != nil {
		return stmt, err
	}
	if _, err := p.Expect(token.LPAREN); err != nil {
		return stmt, err
	}
	if _, err := p.Expect(token.RPAREN); err != nil {
		return stmt, err
	}
	if _, err := p.Expect(token.LBRACE); err != nil {
		return stmt, err
	}
	stmt.Body = p.parseBody()
	if _, err := p.Expect(token.RBRACE); err != nil {
		return stmt, err
	}
	return stmt, nil
}
