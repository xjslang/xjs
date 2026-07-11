package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type IncExpr struct {
	ast.BaseExpr
	Layout struct {
		Increment token.Token
	}
	Left ast.Expr
}

func ParseIncExpr(p *parser.Parser, left ast.Expr) (node *IncExpr, err error) {
	node = &IncExpr{Left: left}
	if node.Layout.Increment, err = p.Expect(token.INCREMENT); err != nil {
		return
	}
	return
}

func PrintIncExpr(pr *printer.Printer, node *IncExpr) error {
	pr.Print(node.Left, node.Layout.Increment)
	return nil
}
