package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var QUESTION_MARK = token.RegisterType("?")

type TernaryExpr struct {
	ast.BaseExpr
	Layout struct {
		QuestionMark token.Token
		Colon        token.Token
	}
	Cond, Then, Else ast.Expr
}

func ParseTernaryExpr(p *parser.Parser, left ast.Expr) (node *TernaryExpr, err error) {
	node = &TernaryExpr{Cond: left}
	if node.Layout.QuestionMark, err = p.Expect(QUESTION_MARK); err != nil {
		return
	}
	if node.Then, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Colon, err = p.Expect(token.COLON); err != nil {
		return
	}
	if node.Else, err = p.ParseExpr(); err != nil {
		return
	}
	return
}

func PrintTernaryExpr(pr *printer.Printer, node *TernaryExpr) error {
	pr.Print(node.Cond)
	pr.Space().Print(node.Layout.QuestionMark)
	pr.Space().Print(node.Then)
	pr.Space().Print(node.Layout.Colon)
	pr.Space().Print(node.Else)
	return nil
}
