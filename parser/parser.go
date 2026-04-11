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
	p.AdvanceToken()
	p.AdvanceToken()
	return p
}

func (p *Parser) AddError(err error) {
	// unimplemented
}

func (p *Parser) AdvanceToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

// Expect checks that the current token matches the expected type,
// advances the position, and returns the new current token.
//
// If the token does not match, it records an error and returns it.
func (p *Parser) Expect(ttype token.TokenType) (token.Token, error) {
	if p.CurrentToken.Type != ttype {
		err := errors.New("Expected " + ttype.String() + ", got " + p.CurrentToken.Type.String())
		p.AddError(err)
		return token.Token{}, err
	}
	p.AdvanceToken()
	return p.CurrentToken, nil
}

func (p *Parser) ParseProgram() *ast.BlockStatement {
	return p.parseBody()
}
