package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type SemiStmt struct {
	ast.BaseStmt
	Layout struct {
		Semi token.Token
	}
}

func ParseSemiStmt(p *parser.Parser) (node *SemiStmt, err error) {
	node = &SemiStmt{}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return node, nil
}

func PrintSemiStmt(p *printer.Printer, node *SemiStmt) error {
	p.Line().Print(node.Layout.Semi)
	return nil
}
