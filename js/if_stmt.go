package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	IF   = token.RegisterType("if")
	ELSE = token.RegisterType("else")
)

type IfStmt struct {
	ast.BaseStmt
	Layout struct {
		If     token.Token
		Lparen token.Token
		Rparen token.Token
		Else   token.Token
	}
	Cond ast.Expr
	Then ast.Stmt
	Else ast.Stmt
}

func ParseIfStmt(p *parser.Parser) (_ *IfStmt, err error) {
	node := &IfStmt{}
	// if
	if node.Layout.If, err = p.Expect(IF); err != nil {
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
	// else
	if p.CurrentToken.Type == ELSE {
		node.Layout.Else = p.CurrentToken
		p.AdvanceToken()
		if node.Else, err = p.ParseStmt(); err != nil {
			return
		}
	}
	return node, nil
}

func PrintIfStmt(p *printer.Printer, node *IfStmt) {
	// if (condition) stmt
	p.LnPrint(node.Layout.If)
	p.SpPrint(node.Layout.Lparen)
	p.Print(node.Cond, node.Layout.Rparen)
	p.SpPrint(node.Then)
	// else
	if node.Else != nil {
		p.SpPrint(node.Layout.Else)
		p.SpPrint(node.Else)
	}
}
