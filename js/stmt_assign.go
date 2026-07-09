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

func ParseAssignStmt(p *parser.Parser) (node *AssignStmt, err error) {
	node = &AssignStmt{}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	if node.Value, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return node, nil
}

func PrintAssignStmt(pr *printer.Printer, node *AssignStmt) error {
	pr.Line().Print(node.Name)
	pr.Space().Print(node.Layout.Assign)
	pr.Space().Print(node.Value)
	pr.Print(node.Layout.Semi)
	return nil
}
