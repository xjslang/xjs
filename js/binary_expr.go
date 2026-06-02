package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type BinaryExpr struct {
	ast.ExprNode
	Tokens struct {
		Op token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseBinaryExpr(p *parser.Parser, left ast.Expr) (_ ast.Expr, err error) {
	op := p.CurrentToken
	node := &BinaryExpr{Left: left}
	node.Tokens.Op = op
	p.AdvanceToken()
	if node.Right, err = ParseRightExpr(p, op.Type.Precedence()); err != nil {
		return
	}
	return node, nil
}

func PrintBinaryExpr(p *printer.Printer, node *BinaryExpr) {
	p.Print(node.Left)
	p.SpPrint(node.Tokens.Op)
	p.SpPrint(node.Right)
}
