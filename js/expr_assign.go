package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type AssignExpr struct {
	ast.BaseExpr
	Layout struct {
		Assign token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseAssignExpr(p *parser.Parser, left ast.Expr) (node *AssignExpr, err error) {
	node = &AssignExpr{Left: left}
	if node.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	if node.Right, err = p.ParseExpr(); err != nil {
		return
	}
	return node, nil
}

func PrintAssignExpr(pr *printer.Printer, node *AssignExpr) error {
	pr.Log("(")
	defer pr.Log(")")
	pr.Print(node.Left).Space().Print(node.Layout.Assign).Space().Print(node.Right)
	return nil
}
