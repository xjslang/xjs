package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type MemberExpr struct {
	ast.BaseExpr
	Layout struct {
		Dot token.Token
	}
	Left  ast.Expr
	Right *Ident
}

func ParseMemberExpr(p *parser.Parser, left ast.Expr) (node *MemberExpr, err error) {
	node = &MemberExpr{Left: left}
	if node.Layout.Dot, err = p.Expect(token.DOT); err != nil {
		return
	}
	if node.Right, err = ParseMemberKey(p); err != nil {
		return
	}
	return
}

func PrintMemberExpr(p *printer.Printer, node *MemberExpr) {
	p.Print(node.Left, node.Layout.Dot, node.Right)
}
