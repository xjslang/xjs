package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type DecExpr struct {
	ast.BaseExpr
	Layout struct {
		Decrement token.Token
	}
	Left ast.Expr
}

func ParseDecExpr(p *parser.Parser, left ast.Expr) (node *DecExpr, err error) {
	node = &DecExpr{Left: left}
	if node.Layout.Decrement, err = p.Expect(token.DECREMENT); err != nil {
		return
	}
	return
}

func PrintDecExpr(p *printer.Printer, node *DecExpr) error {
	p.Print(node.Left).Print(node.Layout.Decrement)
	return nil
}
