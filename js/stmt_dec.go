package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type DecStmt struct {
	ast.BaseStmt
	Layout struct {
		Decrement token.Token
	}
	Name *Ident
}

func ParseDecStmt(p *parser.Parser) (node *DecStmt, err error) {
	node = &DecStmt{}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Decrement, err = p.Expect(token.DECREMENT); err != nil {
		return
	}
	return node, nil
}

func PrintDecStmt(p *printer.Printer, node *DecStmt) {
	p.LnPrint(node.Name)
	p.Print(node.Layout.Decrement)
}
