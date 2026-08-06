package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var ASYNC = token.RegisterType("async")

type PrefixAsyncOp struct {
	ast.BaseExpr
	Layout struct {
		Async token.Token
	}
	Expr ast.Expr
}

func ParsePrefixAsyncOp(p *parser.Parser) (node *PrefixAsyncOp, err error) {
	node = &PrefixAsyncOp{}
	if node.Layout.Async, err = p.Expect(ASYNC); err != nil {
		return
	}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	return
}

func PrintPrefixAsyncOp(pr *printer.Printer, node *PrefixAsyncOp) error {
	pr.Print(node.Layout.Async)
	pr.Space().Print(node.Expr)
	return nil
}
