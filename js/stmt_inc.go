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
	if _, err = ExpectSemi(p); err != nil {
		return
	}
	return
}

func PrintIncStmt(p *printer.Printer, node *IncStmt) {
	p.Line().Print(node.Name)
	p.Print(node.Layout.Increment)
	p.Print(";")
}
