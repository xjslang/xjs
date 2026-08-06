package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type PostfixIncOp struct {
	ast.BaseExpr
	Layout struct {
		Increment token.Token
	}
	Left ast.Expr
}

func ParsePostfixIncOp(p *parser.Parser, left ast.Expr) (node *PostfixIncOp, err error) {
	node = &PostfixIncOp{Left: left}
	if node.Layout.Increment, err = p.Expect(token.INCREMENT); err != nil {
		return
	}
	return
}

func PrintPostfixIncOp(pr *printer.Printer, node *PostfixIncOp) error {
	pr.Print(node.Left, node.Layout.Increment)
	return nil
}
