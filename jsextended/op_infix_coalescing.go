package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var COALESCING = token.RegisterType("COALESCING", "??")

type InfixCoalescingOp struct {
	ast.BaseExpr
	Layout struct {
		Coalescing token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseInfixCoalescingOp(p *parser.Parser, left ast.Expr) (node *InfixCoalescingOp, err error) {
	node = &InfixCoalescingOp{Left: left}
	if node.Layout.Coalescing, err = p.Expect(COALESCING); err != nil {
		return
	}
	if node.Right, err = js.ParseRightExpr(p, COALESCING.Precedence()); err != nil {
		return
	}
	return
}

func PrintInfixCoalescingOp(pr *printer.Printer, node *InfixCoalescingOp) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Layout.Coalescing)
	pr.Space().Print(node.Right)
	return nil
}
