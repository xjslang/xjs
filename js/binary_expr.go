package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type BinaryExpr struct {
	ast.ExprNode
	LeftValue  ast.Expr
	Operator   token.Token
	RightValue ast.Expr
}

func (node *BinaryExpr) Type() string {
	return "BinaryExpr"
}

func ParseBinaryExpr(p *parser.Parser, leftVal ast.Expr) (node ast.Expr, err error) {
	op := p.CurrentToken
	nodeExpr := &BinaryExpr{
		LeftValue: leftVal,
		Operator:  op,
	}
	p.AdvanceToken()
	if nodeExpr.RightValue, err = ParseRightExpr(p, op.Type.Precedence()); err != nil {
		return
	}
	node = nodeExpr
	return
}

func PrintBinaryExpr(p *printer.Printer, node *BinaryExpr) {
	p.Print(node.LeftValue)
	p.SpPrint(node.Operator)
	p.SpPrint(node.RightValue)
}
