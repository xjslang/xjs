package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type BinaryExpr struct {
	ast.BaseExpr
	Left  ast.Expr
	Op    token.Token
	Right ast.Expr
}

func ParseBinaryExpr(p *parser.Parser, left ast.Expr) (_ *BinaryExpr, err error) {
	op := p.CurrentToken
	node := &BinaryExpr{Left: left, Op: op}
	p.AdvanceToken()
	if node.Right, err = ParseRightExpr(p, op.Type.Precedence()); err != nil {
		return
	}
	return node, nil
}

func PrintBinaryExpr(p *printer.Printer, node *BinaryExpr) {
	switch node.Op.Type {
	case token.DOT:
		p.Print(node.Left, node.Op, node.Right)
	default:
		p.Print(node.Left)
		p.SpPrint(node.Op)
		p.SpPrint(node.Right)
	}
}
