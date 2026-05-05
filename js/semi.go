package js

import (
	"errors"

	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func ExpectSemi(p *parser.Parser) error {
	if advanceSemi(p) {
		return nil
	}
	msg := "Expected statement terminator"
	p.AddError(msg)
	return errors.New(msg)
}

func AdvanceToStatementEnd(p *parser.Parser) {
	for !advanceSemi(p) {
		p.AdvanceToken()
	}
}

func advanceSemi(p *parser.Parser) bool {
	if p.CurrentToken.Type == token.SEMICOLON {
		p.AdvanceToken()
		return true
	}
	if p.CurrentToken.Type == token.EOF || p.CurrentToken.AfterNewline {
		return true
	}
	if p.InScope(parser.BlockScope) && p.CurrentToken.Type == token.RBRACE {
		return true
	}
	return false
}
