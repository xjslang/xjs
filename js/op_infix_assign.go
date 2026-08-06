package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type InfixAssignOp struct {
	ast.BaseExpr
	Layout struct {
		Assign token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseInfixAssignOp(p *parser.Parser, left ast.Expr) (node *InfixAssignOp, err error) {
	node = &InfixAssignOp{Left: left}
	if node.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	if node.Right, err = p.ParseExpr(); err != nil {
		return
	}
	return node, nil
}

func PrintInfixAssignOp(pr *printer.Printer, node *InfixAssignOp) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Layout.Assign)
	pr.Space().Print(node.Right)
	return nil
}
