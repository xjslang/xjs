package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var DELETE = token.RegisterType("delete")

type DeleteExpr struct {
	ast.BaseExpr
	Layout struct {
		Delete token.Token
	}
	Value ast.Expr
}

func ParseDeleteExpr(p *parser.Parser) (node *DeleteExpr, err error) {
	node = &DeleteExpr{}
	if node.Layout.Delete, err = p.Expect(DELETE); err != nil {
		return
	}
	if node.Value, err = ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintDeleteExpr(pr *printer.Printer, node *DeleteExpr) error {
	pr.Print(node.Layout.Delete)
	pr.Space().Print(node.Value)
	return nil
}
