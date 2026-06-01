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
	IfToken     token.Token
	LparenToken token.Token
	RparenToken token.Token
	ElseToken   token.Token

	CondExpr ast.Node
	ThenStmt ast.Node
	ElseStmt ast.Node
}

func (node *IfStmt) Type() string {
	return "IfStmt"
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
	// else
	if p.CurrentToken.Type == ELSE {
		node.ElseToken = p.CurrentToken
		p.AdvanceToken()
		if node.ElseStmt, err = p.ParseStmt(); err != nil {
			return
		}
	}
	return node, nil
}

func PrintIfStmt(p *printer.Printer, node *IfStmt) {
	// if (condition) stmt
	p.LnPrint(node.IfToken)
	p.SpPrint(node.LparenToken)
	p.Print(node.CondExpr, node.RparenToken)
	p.SpPrint(node.ThenStmt)
	// else
	if node.ElseStmt != nil {
		p.SpPrint(node.ElseToken)
		p.SpPrint(node.ElseStmt)
	}
}
