package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
)

type ExprStmt struct {
	ast.BaseStmt
	Expr ast.Expr
}

func ParseExprStmt(p *parser.Parser) (_ *ExprStmt, err error) {
	node := &ExprStmt{}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	return node, nil
}

func PrintExprStmt(p *printer.Printer, node *ExprStmt) {
	p.LnPrint(node.Expr)
}
