package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var EXPO = token.RegisterType("EXPO", "**")

type InfixExpoOp struct {
	ast.BaseExpr
	Layout struct {
		Expo token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseInfixExpoOp(p *parser.Parser, left ast.Expr) (node *InfixExpoOp, err error) {
	node = &InfixExpoOp{Left: left}
	if node.Layout.Expo, err = p.Expect(EXPO); err != nil {
		return
	}
	if node.Right, err = js.ParseRightExpr(p, EXPO.Precedence()); err != nil {
		return
	}
	return
}

func PrintInfixExpoOp(pr *printer.Printer, node *InfixExpoOp) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Layout.Expo)
	pr.Space().Print(node.Right)
	return nil
}
