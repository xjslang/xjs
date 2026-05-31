package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type UnaryExpr struct {
	Operator token.Token
	Value    ast.Node
}

func (node *UnaryExpr) Type() string {
	return "UnaryExpr"
}

func ParseUnaryExpr(p *parser.Parser) (node ast.Node, err error) {
	nodeExpr := &UnaryExpr{Operator: p.CurrentToken}
	p.AdvanceToken()
	if nodeExpr.Value, err = ParseValue(p); err != nil {
		return
	}
	node = nodeExpr
	return
}

func PrintUnaryExpr(p *printer.Printer, node *UnaryExpr) {
	p.Print(node.Operator, node.Value)
}
