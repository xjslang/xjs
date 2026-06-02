package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var LET = token.RegisterType("let")

type LetStmt struct {
	ast.BaseStmt
	Tokens struct {
		Let    token.Token
		Assign token.Token
		Semi   token.Token
	}
	Name  *Ident
	Value ast.Expr
}

func ParseLetStmt(p *parser.Parser) (_ *LetStmt, err error) {
	node := &LetStmt{}
	if node.Tokens.Let, err = p.Expect(LET); err != nil {
		return
	}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Tokens.Assign, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	if node.Tokens.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintLetStmt(p *printer.Printer, node *LetStmt) {
	p.LnPrint(node.Tokens.Let)
	p.SpPrint(node.Name)
	p.SpPrint(node.Tokens.Assign)
	p.SpPrint(node.Value)
	p.Print(node.Tokens.Semi)
}
