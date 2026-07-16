package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var SPREAD = token.RegisterType("...")

type SpreadExpr struct {
	ast.BaseExpr
	Layout struct {
		Spread token.Token
	}
	Value ast.Expr
}

func ParseSpreadExpr(p *parser.Parser) (node *SpreadExpr, err error) {
	node = &SpreadExpr{}
	if node.Layout.Spread, err = p.Expect(SPREAD); err != nil {
		return
	}
	if node.Value, err = js.ParseValue(p); err != nil {
		return
	}
	return
}

func PrintSpreadExpr(pr *printer.Printer, node *SpreadExpr) error {
	pr.Print(node.Layout.Spread, node.Value)
	return nil
}
