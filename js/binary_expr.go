package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type BinaryExpr struct {
	LeftValue  ast.Node
	Operator   token.Token
	RightValue ast.Node
}

func (node *BinaryExpr) Type() string {
	return "BinaryExpr"
}

func ParseInfixExpr(p *parser.Parser, leftVal ast.Node) (node ast.Node, err error) {
	op := p.CurrentToken
	nodeExpr := &BinaryExpr{
		LeftValue: leftVal,
		Operator:  op,
	}
	p.AdvanceToken()
	if nodeExpr.RightValue, err = ParseRightValue(p, op.Type.Precedence()); err != nil {
		return
	}
	node = nodeExpr
	return
}

func PrintBinaryExpr(p *printer.Printer, node *BinaryExpr) {
	p.Print(node.LeftValue)
	p.SpPrint(node.Operator, node.RightValue)
}
