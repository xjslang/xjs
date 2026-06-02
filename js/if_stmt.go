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
	Tokens struct {
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
	if node.Tokens.If, err = p.Expect(IF); err != nil {
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
	// else
	if p.CurrentToken.Type == ELSE {
		node.Tokens.Else = p.CurrentToken
		p.AdvanceToken()
		if node.Else, err = p.ParseStmt(); err != nil {
			return
		}
	}
	return node, nil
}

func PrintIfStmt(p *printer.Printer, node *IfStmt) {
	// if (condition) stmt
	p.LnPrint(node.Tokens.If)
	p.SpPrint(node.Tokens.Lparen)
	p.Print(node.Cond, node.Tokens.Rparen)
	p.SpPrint(node.Then)
	// else
	if node.Else != nil {
		p.SpPrint(node.Tokens.Else)
		p.SpPrint(node.Else)
	}
}
