package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type BinaryExpr struct {
	ast.ExprNode
	Left  ast.Expr
	Op    token.Token
	Right ast.Expr
}

func (node *BinaryExpr) Type() string {
	return "BinaryExpr"
}

func ParseBinaryExpr(p *parser.Parser, left ast.Expr) (node ast.Expr, err error) {
	op := p.CurrentToken
	nodeExpr := &BinaryExpr{
		Left: left,
		Op:   op,
	}
	p.AdvanceToken()
	if nodeExpr.Right, err = ParseRightExpr(p, op.Type.Precedence()); err != nil {
		return
	}
	node = nodeExpr
	return
}

func PrintBinaryExpr(p *printer.Printer, node *BinaryExpr) {
	p.Print(node.Left)
	p.SpPrint(node.Op)
	p.SpPrint(node.Right)
}
