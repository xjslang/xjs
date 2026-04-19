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
	// call twice to update CurrentToken and PeekToken
	p.advanceToken()
	p.advanceToken()
	return p
}

func (p *Parser) addError(err error) {
	// TODO: implement
}

func (p *Parser) advanceToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

// Expect checks that the current token matches the expected type,
// advances the position, and returns the token.
//
// If the token does not match, it records an error and returns it.
func (p *Parser) expect(ttype token.TokenType) (token.Token, error) {
	if p.CurrentToken.Type != ttype {
		err := errors.New("Expected " + ttype.String() + ", got " + p.CurrentToken.Type.String())
		p.addError(err)
		return token.Token{}, err
	}
	tok := p.CurrentToken
	p.advanceToken()
	return tok, nil
}

func (p *Parser) expectASI() error {
	if p.CurrentToken.Type == token.SEMICOLON {
		p.advanceToken()
		return nil
	}
	if p.CurrentToken.Type == token.EOF || p.CurrentToken.AfterNewline {
		return nil
	}
	err := errors.New("Expected semicolon, newline, or EOF, got " + p.CurrentToken.Type.String())
	p.addError(err)
	return err
}

func (p *Parser) ParseProgram() *ast.BlockStatement {
	return p.parseBody()
}

func (p *Parser) ParseExpression() ast.Expression {
	switch p.CurrentToken.Type {
	case token.NUMBER:
		val := p.CurrentToken.Literal
		p.advanceToken()
		return &ast.IntegerLiteral{Value: val}
	case token.STRING:
		val := p.CurrentToken.Literal
		p.advanceToken()
		return &ast.StringLiteral{Value: val}
	}
	return nil
}

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
	p.advanceToken() // consume token.LET
	ident, err := p.expect(token.IDENT)
	if err != nil {
		return stmt, err
	}
	stmt.Name = &ast.Identifier{Value: ident.Literal}
	if _, err := p.expect(token.ASSIGN); err != nil {
		return stmt, err
	}
	stmt.Value = p.ParseExpression()
	if err := p.expectASI(); err != nil {
		return stmt, err
	}
	return stmt, nil
}

func (p *Parser) parseFunction() (*ast.FunctionDeclaration, error) {
	stmt := &ast.FunctionDeclaration{}
	p.advanceToken() // consume token.FUNCTION
	ident, err := p.expect(token.IDENT)
	if err != nil {
		return stmt, err
	}
	stmt.Name = &ast.Identifier{Value: ident.Literal}
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
