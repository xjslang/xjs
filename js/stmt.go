package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Stmt struct {
	SemiToken token.Token

	Expr ast.Node
}

func (node *Stmt) Type() string {
	return "Stmt"
}

func ParseStmt(p *parser.Parser) (_ *Stmt, err error) {
	node := &Stmt{}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintStmt(p *printer.Printer, node *Stmt) {
	p.LnPrint(node.Expr)
	p.Print(node.SemiToken)
}
