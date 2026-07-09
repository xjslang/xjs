package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var WHILE = token.RegisterType("while")

type WhileStmt struct {
	ast.BaseStmt
	Layout struct {
		While  token.Token
		Lparen token.Token
		Rparen token.Token
	}
	Cond ast.Expr
	Then ast.Stmt
}

func ParseWhileStmt(p *parser.Parser) (node *WhileStmt, err error) {
	node = &WhileStmt{}
	// while
	if node.Layout.While, err = p.Expect(WHILE); err != nil {
		return
	}
	// (condition)
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Cond, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	// then
	if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	return node, nil
}

func PrintWhileStmt(pr *printer.Printer, node *WhileStmt) error {
	pr.Line().Print(node.Layout.While)
	pr.Space().Print(node.Layout.Lparen)
	pr.Print(node.Cond, node.Layout.Rparen)
	pr.Space().Print(node.Then)
	return nil
}
