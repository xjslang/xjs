package js

import (
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func AdvanceToStmtEnd(p *parser.Parser) {
	for {
		typ := p.CurrentToken.Type
		if typ == token.SEMICOLON {
			p.AdvanceToken()
			break
		}
		if typ == token.EOF || typ == token.RBRACE || typ == token.LBRACE || p.CurrentToken.AfterNewline {
			break
		}
		p.AdvanceToken()
	}
}
