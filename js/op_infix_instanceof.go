package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var INSTANCEOF = token.RegisterType("INSTANCEOF", "instanceof")

type InfixInstanceofOp struct {
	ast.BaseExpr
	Layout struct {
		Instanceof token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseInfixInstanceofOp(p *parser.Parser, left ast.Expr) (node *InfixInstanceofOp, err error) {
	node = &InfixInstanceofOp{Left: left}
	if node.Layout.Instanceof, err = p.Expect(INSTANCEOF); err != nil {
		return
	}
	if node.Right, err = ParseRightExpr(p, INSTANCEOF.Precedence()); err != nil {
		return
	}
	return
}

func PrintInfixInstanceofOp(pr *printer.Printer, node *InfixInstanceofOp) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Layout.Instanceof)
	pr.Space().Print(node.Right)
	return nil
}
