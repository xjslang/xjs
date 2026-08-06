package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type InfixOp struct {
	ast.BaseExpr
	Left  ast.Expr
	Op    token.Token
	Right ast.Expr
}

func ParseInfixOp(p *parser.Parser, left ast.Expr) (node *InfixOp, err error) {
	op := p.CurrentToken
	node = &InfixOp{Left: left, Op: op}
	p.AdvanceToken()
	if node.Right, err = ParseRightExpr(p, op.Type.Precedence()); err != nil {
		return
	}
	return node, nil
}

func PrintInfixOp(pr *printer.Printer, node *InfixOp) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Op)
	pr.Space().Print(node.Right)
	return nil
}
