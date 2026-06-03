package js

import (
	"errors"

	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

var forHeaderScope = parser.RegisterScope()

func ExpectSemi(p *parser.Parser) (token.Token, error) {
	tok := p.CurrentToken
	if tok.Type == token.SEMICOLON {
		if !p.InScope(forHeaderScope) {
			p.AdvanceToken()
		}
		return tok, nil
	}
	if tok.Type == token.EOF || tok.AfterNewline ||
		p.InScope(blockScope) && tok.Type == token.RBRACE ||
		p.InScope(forHeaderScope) && tok.Type == token.RPAREN {
		tok = token.Token{Type: token.SEMICOLON, Literal: token.SEMICOLON.String(), Position: tok.Position}
		return tok, nil
	}
	msg := "Expected statement terminator"
	p.AddError(msg)
	return tok, errors.New(msg)
}
