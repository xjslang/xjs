package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type UnaryExpr struct {
	ast.BaseExpr
	Op    token.Token
	Value ast.Expr
}

func ParseUnaryExpr(p *parser.Parser) (node *UnaryExpr, err error) {
	node = &UnaryExpr{}
	node.Op = p.CurrentToken
	p.AdvanceToken()
	if node.Value, err = ParseValue(p); err != nil {
		return
	}
	return node, nil
}

func PrintUnaryExpr(p *printer.Printer, node *UnaryExpr) error {
	p.Print(node.Op, node.Value)
	return nil
}
