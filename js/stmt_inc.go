package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type IncStmt struct {
	ast.BaseStmt
	Layout struct {
		Increment token.Token
	}
	Name *Ident
}

func ParseIncStmt(p *parser.Parser) (node *IncStmt, err error) {
	node = &IncStmt{}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Increment, err = p.Expect(token.INCREMENT); err != nil {
		return
	}
	return node, nil
}

func PrintIncStmt(p *printer.Printer, node *IncStmt) {
	p.LnPrint(node.Name)
	p.Print(node.Layout.Increment)
}
