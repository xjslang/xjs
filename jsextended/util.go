package jsextended

import (
	"slices"

	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// token types allowed to be used as identifiers
var allowedIdTypes = []token.Type{js.LET, YIELD}

func ParseIdent(p *parser.Parser) (node *js.Ident, err error) {
	if slices.Contains(allowedIdTypes, p.CurrentToken.Type) {
		node = &js.Ident{Token: p.CurrentToken}
		p.AdvanceToken()
		return
	}
	return js.ParseIdent(p)
}
