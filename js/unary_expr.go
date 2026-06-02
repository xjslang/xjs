package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type UnaryExpr struct {
	ast.ExprNode
	Op    token.Token
	Value ast.Expr
}

func (node *UnaryExpr) Type() string {
	return "UnaryExpr"
}

func ParseUnaryExpr(p *parser.Parser) (node ast.Expr, err error) {
	nodeExpr := &UnaryExpr{Op: p.CurrentToken}
	p.AdvanceToken()
	if nodeExpr.Value, err = ParseValue(p); err != nil {
		return
	}
	node = nodeExpr
	return
}

func PrintUnaryExpr(p *printer.Printer, node *UnaryExpr) {
	p.Print(node.Op, node.Value)
}
