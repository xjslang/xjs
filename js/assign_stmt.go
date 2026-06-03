package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type AssignStmt struct {
	ast.BaseStmt
	Layout struct {
		Assign token.Token
		Semi   token.Token
	}
	Name  *Ident
	Value ast.Expr
}

func ParseAssignStmt(p *parser.Parser) (_ *AssignStmt, err error) {
	node := &AssignStmt{}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	if node.Value, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintAssignStmt(p *printer.Printer, node *AssignStmt) {
	p.LnPrint(node.Name)
	p.SpPrint(node.Layout.Assign)
	p.SpPrint(node.Value)
	p.Print(node.Layout.Semi)
}
