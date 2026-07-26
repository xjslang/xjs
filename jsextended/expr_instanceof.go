package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var INSTANCEOF = token.RegisterType("instanceof")

type InstanceofExpr struct {
	ast.BaseExpr
	Layout struct {
		Instanceof token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseInstanceofExpr(p *parser.Parser, left ast.Expr) (node *InstanceofExpr, err error) {
	node = &InstanceofExpr{Left: left}
	if node.Layout.Instanceof, err = p.Expect(INSTANCEOF); err != nil {
		return
	}
	if node.Right, err = js.ParseRightExpr(p, INSTANCEOF.Precedence()); err != nil {
		return
	}
	return
}

func PrintInstanceofExpr(pr *printer.Printer, node *InstanceofExpr) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Layout.Instanceof)
	pr.Space().Print(node.Right)
	return nil
}
