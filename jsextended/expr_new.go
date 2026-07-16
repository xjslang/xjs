package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var NEW = token.RegisterType("new")

type NewExpr struct {
	ast.BaseExpr
	Layout struct {
		New token.Token
	}
	Value ast.Expr
}

func ParseNewExpr(p *parser.Parser) (node *NewExpr, err error) {
	node = &NewExpr{}
	if node.Layout.New, err = p.Expect(NEW); err != nil {
		return
	}
	if node.Value, err = js.ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintNewExpr(pr *printer.Printer, node *NewExpr) error {
	pr.Print(node.Layout.New)
	pr.Space().Print(node.Value)
	return nil
}
