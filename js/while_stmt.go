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
	Tokens struct {
		While  token.Token
		Lparen token.Token
		Rparen token.Token
	}
	Cond ast.Expr
	Then ast.Stmt
}

func ParseWhileStmt(p *parser.Parser) (_ *WhileStmt, err error) {
	node := &WhileStmt{}
	// while
	if node.Tokens.While, err = p.Expect(WHILE); err != nil {
		return
	}
	// (condition)
	if node.Tokens.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Cond, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Tokens.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	// then
	if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	return node, nil
}

func PrintWhileStmt(p *printer.Printer, node *WhileStmt) {
	p.LnPrint(node.Tokens.While)
	p.SpPrint(node.Tokens.Lparen)
	p.Print(node.Cond, node.Tokens.Rparen)
	p.SpPrint(node.Then)
}
