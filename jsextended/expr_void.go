package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var VOID = token.RegisterType("void")

type VoidExpr struct {
	ast.BaseExpr
	Layout struct {
		Void token.Token
	}
	Expr ast.Expr
}

func ParseVoidExpr(p *parser.Parser) (node *VoidExpr, err error) {
	node = &VoidExpr{}
	if node.Layout.Void, err = p.Expect(VOID); err != nil {
		return
	}
	if node.Expr, err = js.ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintVoidExpr(pr *printer.Printer, node *VoidExpr) error {
	pr.Print(node.Layout.Void)
	pr.Space().Print(node.Expr)
	return nil
}
