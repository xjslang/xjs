package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
)

func ParseExpr(p *parser.Parser) (val ast.Node, err error) {
	if val, err = ParseValue(p); err != nil {
		return
	}
	typ := p.CurrentToken.Type
	for typ.IsInfixOp() {
		if val, err = p.ParseInfixExpr(val); err != nil {
			return
		}
		typ = p.CurrentToken.Type
	}
	return
}
