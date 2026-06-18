package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

type SelfClosingStmt interface {
	ast.Stmt
	SelfClosing() bool
}

func ParseStmt(p *parser.Parser) (ast.Stmt, error) {
	if p.CurrentToken.Type == token.LBRACE {
		return ParseBlockStmt(p)
	}
	return ParseExprStmt(p)
}
