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
	ast.StmtNode
	IfToken     token.Token
	LparenToken token.Token
	RparenToken token.Token
	ElseToken   token.Token

	Cond ast.Expr
	Then ast.Stmt
	Else ast.Stmt
}

func ParseIfStmt(p *parser.Parser) (_ *IfStmt, err error) {
	node := &IfStmt{}
	// if
	if node.IfToken, err = p.Expect(IF); err != nil {
		return
	}
	// (condition)
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Cond, err = p.ParseExpr(); err != nil {
		return
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	// then
	if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	// else
	if p.CurrentToken.Type == ELSE {
		node.ElseToken = p.CurrentToken
		p.AdvanceToken()
		if node.Else, err = p.ParseStmt(); err != nil {
			return
		}
	}
	return node, nil
}

func PrintIfStmt(p *printer.Printer, node *IfStmt) {
	// if (condition) stmt
	p.LnPrint(node.IfToken)
	p.SpPrint(node.LparenToken)
	p.Print(node.Cond, node.RparenToken)
	p.SpPrint(node.Then)
	// else
	if node.Else != nil {
		p.SpPrint(node.ElseToken)
		p.SpPrint(node.Else)
	}
}
