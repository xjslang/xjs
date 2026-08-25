package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var SPREAD = token.RegisterType("SPREAD", "...")

type PrefixSpreadOp struct {
	ast.BaseExpr
	Layout struct {
		Spread token.Token
	}
	Value ast.Expr
}

func ParsePrefixSpreadOp(p *parser.Parser) (node *PrefixSpreadOp, err error) {
	node = &PrefixSpreadOp{}
	if node.Layout.Spread, err = p.Expect(SPREAD); err != nil {
		return
	}
	if node.Value, err = ParseValue(p); err != nil {
		return
	}
	return
}

func PrintPrefixSpreadOp(pr *printer.Printer, node *PrefixSpreadOp) error {
	pr.Print(node.Layout.Spread, node.Value)
	return nil
}
