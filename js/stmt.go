package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ExprStmt struct {
	SemiToken token.Token

	Expr ast.Node
}

func (node *ExprStmt) Type() string {
	return "ExprStmt"
}

func ParseExprStmt(p *parser.Parser) (_ *ExprStmt, err error) {
	node := &ExprStmt{}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintExprStmt(p *printer.Printer, node *ExprStmt) {
	p.LnPrint(node.Expr)
	p.Print(node.SemiToken)
}
