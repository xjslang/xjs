package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ExprStmt struct {
	ast.BaseStmt
	Layout struct {
		Semi token.Token
	}
	Expr ast.Expr
}

func ParseExprStmt(p *parser.Parser) (node *ExprStmt, err error) {
	node = &ExprStmt{}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return
}

func PrintExprStmt(p *printer.Printer, node *ExprStmt) error {
	p.Line().Print(node.Expr)
	p.Print(node.Layout.Semi)
	return nil
}
