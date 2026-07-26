package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var COALESCING = token.RegisterType("??")

type CoalescingExpr struct {
	ast.BaseExpr
	Layout struct {
		Coalescing token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseCoalescingExpr(p *parser.Parser, left ast.Expr) (node *CoalescingExpr, err error) {
	node = &CoalescingExpr{Left: left}
	if node.Layout.Coalescing, err = p.Expect(COALESCING); err != nil {
		return
	}
	if node.Right, err = js.ParseRightExpr(p, COALESCING.Precedence()); err != nil {
		return
	}
	return
}

func PrintCoalescingExpr(pr *printer.Printer, node *CoalescingExpr) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Layout.Coalescing)
	pr.Space().Print(node.Right)
	return nil
}
