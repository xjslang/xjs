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
	Layout struct {
		Let    token.Token
		Assign token.Token
	}
	Name  *Ident
	Value ast.Expr
}

func ParseLetStmt(p *parser.Parser) (node *LetStmt, err error) {
	node = &LetStmt{}
	if node.Layout.Let, err = p.Expect(LET); err != nil {
		return
	}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	return node, nil
}

func PrintLetStmt(p *printer.Printer, node *LetStmt) {
	p.Line().Print(node.Layout.Let)
	p.Space().Print(node.Name)
	p.Space().Print(node.Layout.Assign)
	p.Space().Print(node.Value)
}
