package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var WHILE = token.RegisterType("while")

type WhileStmt struct {
	ast.StmtNode
	WhileToken  token.Token
	LparenToken token.Token
	RparenToken token.Token

	CondExpr ast.Expr
	ThenStmt ast.Stmt
}

func (node *WhileStmt) Type() string {
	return "WhileStmt"
}

func ParseWhileStmt(p *parser.Parser) (_ *WhileStmt, err error) {
	node := &WhileStmt{}
	// while
	if node.WhileToken, err = p.Expect(WHILE); err != nil {
		return
	}
	// (condition)
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.CondExpr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	// then
	if node.ThenStmt, err = p.ParseStmt(); err != nil {
		return
	}
	return node, nil
}

func PrintWhileStmt(p *printer.Printer, node *WhileStmt) {
	p.LnPrint(node.WhileToken)
	p.SpPrint(node.LparenToken)
	p.Print(node.CondExpr, node.RparenToken)
	p.SpPrint(node.ThenStmt)
}
