package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var VOID = token.RegisterType("VOID", "void")

type PrefixVoidOp struct {
	ast.BaseExpr
	Layout struct {
		Void token.Token
	}
	Expr ast.Expr
}

func ParsePrefixVoidOp(p *parser.Parser) (node *PrefixVoidOp, err error) {
	node = &PrefixVoidOp{}
	if node.Layout.Void, err = p.Expect(VOID); err != nil {
		return
	}
	if node.Expr, err = ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintPrefixVoidOp(pr *printer.Printer, node *PrefixVoidOp) error {
	pr.Print(node.Layout.Void)
	pr.Space().Print(node.Expr)
	return nil
}
