package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type PrefixOp struct {
	ast.BaseExpr
	Op    token.Token
	Value ast.Expr
}

func ParsePrefixOp(p *parser.Parser) (node *PrefixOp, err error) {
	node = &PrefixOp{}
	node.Op = p.CurrentToken
	p.AdvanceToken()
	if node.Value, err = ParseValue(p); err != nil {
		return
	}
	return node, nil
}

func PrintPrefixOp(pr *printer.Printer, node *PrefixOp) error {
	pr.Print(node.Op, node.Value)
	return nil
}
