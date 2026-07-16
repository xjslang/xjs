package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var THROW = token.RegisterType("throw")

type ThrowStmt struct {
	ast.BaseStmt
	Layout struct {
		Throw token.Token
		Semi  token.Token
	}
	Expr ast.Expr
}

func ParseThrowStmt(p *parser.Parser) (node *ThrowStmt, err error) {
	node = &ThrowStmt{}
	if node.Layout.Throw, err = p.Expect(THROW); err != nil {
		return
	}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Semi, err = js.ExpectSemi(p); err != nil {
		return
	}
	return
}

func PrintThrowStmt(pr *printer.Printer, node *ThrowStmt) error {
	pr.Line().Print(node.Layout.Throw)
	pr.Space().Print(node.Expr, node.Layout.Semi)
	return nil
}
