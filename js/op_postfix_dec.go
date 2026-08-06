package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type PostfixDecOp struct {
	ast.BaseExpr
	Layout struct {
		Decrement token.Token
	}
	Left ast.Expr
}

func ParsePostfixDecOp(p *parser.Parser, left ast.Expr) (node *PostfixDecOp, err error) {
	node = &PostfixDecOp{Left: left}
	if node.Layout.Decrement, err = p.Expect(token.DECREMENT); err != nil {
		return
	}
	return
}

func PrintPostfixDecOp(pr *printer.Printer, node *PostfixDecOp) error {
	pr.Print(node.Left, node.Layout.Decrement)
	return nil
}
