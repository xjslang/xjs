package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type InfixDotOp struct {
	ast.BaseExpr
	Layout struct {
		Dot token.Token
	}
	Left  ast.Expr
	Right *Ident
}

func ParseInfixDotOp(p *parser.Parser, left ast.Expr) (node *InfixDotOp, err error) {
	node = &InfixDotOp{Left: left}
	if node.Layout.Dot, err = p.Expect(token.DOT); err != nil {
		return
	}
	if node.Right, err = ParseObjKey(p); err != nil {
		return
	}
	return
}

func PrintInfixDotOp(pr *printer.Printer, node *InfixDotOp) error {
	if v, ok := node.Left.(*Literal); ok && v.Value.Type == token.NUMBER {
		pr.Print("(", v, ")")
	} else {
		pr.Print(node.Left)
	}
	pr.Print(node.Layout.Dot, node.Right)
	return nil
}
