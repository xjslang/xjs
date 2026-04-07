package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	CurrentToken token.Token
	PeekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.NextToken()
	return p
}

func (p *Parser) AddError(err error) {
	// unimplemented
}

func (p *Parser) NextToken() token.Token {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
	return p.CurrentToken
}

func (p *Parser) ParseProgram() *ast.BlockStatement {
	return p.parseBody()
}

func (p *Parser) expect(ttype token.TokenType) (token.Token, error) {
	if p.PeekToken.Type != ttype {
		err := errors.New("Expected " + ttype.String() + ", got " + p.PeekToken.Type.String())
		p.AddError(err)
		return token.Token{}, err
	}
	return p.NextToken(), nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.PeekToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.FUNCTION:
		return p.parseFunction()
	}
	return nil, errors.New("Unexpected token " + p.PeekToken.Literal)
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{}
	p.NextToken() // consume token.LET
	ident, err := p.expect(token.IDENT)
	if err != nil {
		return stmt, err
	}
	stmt.Name = &ast.Identifier{Value: ident.Literal}
	if _, err := p.expect(token.ASSIGN); err != nil {
		return stmt, err
	}
	stmt.Value = p.ParseExpression()
	if _, err := p.expect(token.SEMI); err != nil {
		return stmt, err
	}
	return stmt, nil
}

func (p *Parser) ParseExpression() ast.Expression {
	switch p.PeekToken.Type {
	case token.NUMBER:
		val := p.NextToken().Literal
		return &ast.IntegerLiteral{Value: val}
	case token.STRING:
		val := p.NextToken().Literal
		return &ast.StringLiteral{Value: val}
	}
	return nil
}

func (p *Parser) parseFunction() (*ast.FunctionStatement, error) {
	stmt := &ast.FunctionStatement{}
	p.NextToken() // consume token.FUNCTION
	if _, err := p.expect(token.IDENT); err != nil {
		return stmt, err
	}
	if _, err := p.expect(token.LPAREN); err != nil {
		return stmt, err
	}
	if _, err := p.expect(token.RPAREN); err != nil {
		return stmt, err
	}
	if _, err := p.expect(token.LBRACE); err != nil {
		return stmt, err
	}
	stmt.Body = p.parseBody()
	if _, err := p.expect(token.RBRACE); err != nil {
		return stmt, err
	}
	return stmt, nil
}

func (p *Parser) parseBody() *ast.BlockStatement {
	bodyStmt := &ast.BlockStatement{}
	for stmt, err := p.parseStatement(); err == nil; stmt, err = p.parseStatement() {
		bodyStmt.Statements = append(bodyStmt.Statements, stmt)
	}
	return bodyStmt
}
