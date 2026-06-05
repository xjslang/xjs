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

func ParseSemiStmt(p *parser.Parser) (_ *SemiStmt, err error) {
	node := &SemiStmt{}
	if node.Layout.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintSemiStmt(p *printer.Printer, node *SemiStmt) {
	p.Print(node.Layout.Semi)
}
