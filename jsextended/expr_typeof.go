package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var TYPEOF = token.RegisterType("typeof")

type TypeofExpr struct {
	ast.BaseExpr
	Layout struct {
		Typeof token.Token
	}
	Value ast.Expr
}

func ParseTypeofExpr(p *parser.Parser) (node *TypeofExpr, err error) {
	node = &TypeofExpr{}
	if node.Layout.Typeof, err = p.Expect(TYPEOF); err != nil {
		return
	}
	if node.Value, err = js.ParseValue(p); err != nil {
		return
	}
	return
}

func PrintTypeofExpr(pr *printer.Printer, node *TypeofExpr) error {
	pr.Log("(")
	defer pr.Log(")")
	pr.Print(node.Layout.Typeof)
	pr.Space().Print(node.Value)
	return nil
}
