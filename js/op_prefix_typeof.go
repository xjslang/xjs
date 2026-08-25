package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var TYPEOF = token.RegisterType("TYPEOF", "typeof")

type PrefixTypeofOp struct {
	ast.BaseExpr
	Layout struct {
		Typeof token.Token
	}
	Value ast.Expr
}

func ParsePrefixTypeofOp(p *parser.Parser) (node *PrefixTypeofOp, err error) {
	node = &PrefixTypeofOp{}
	if node.Layout.Typeof, err = p.Expect(TYPEOF); err != nil {
		return
	}
	if node.Value, err = ParseValue(p); err != nil {
		return
	}
	return
}

func PrintPrefixTypeofOp(pr *printer.Printer, node *PrefixTypeofOp) error {
	pr.Print(node.Layout.Typeof)
	pr.Space().Print(node.Value)
	return nil
}
