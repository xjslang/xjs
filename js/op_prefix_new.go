package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var NEW = token.RegisterType("NEW", "new")

type PrefixNewOp struct {
	ast.BaseExpr
	Layout struct {
		New token.Token
	}
	Value ast.Expr
}

func ParsePrefixNewOp(p *parser.Parser) (node *PrefixNewOp, err error) {
	node = &PrefixNewOp{}
	if node.Layout.New, err = p.Expect(NEW); err != nil {
		return
	}
	if node.Value, err = ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintPrefixNewOp(pr *printer.Printer, node *PrefixNewOp) error {
	pr.Print(node.Layout.New)
	pr.Space().Print(node.Value)
	return nil
}
