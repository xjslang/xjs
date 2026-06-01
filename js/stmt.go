package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func ParseStmt(p *parser.Parser) (_ ast.Node, err error) {
	if p.CurrentToken.Type == token.LBRACE {
		return ParseBlockStmt(p)
	}
	return ParseExprStmt(p)
}
