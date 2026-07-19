package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var AWAIT = token.RegisterType("await")

type AwaitExpr struct {
	ast.BaseExpr
	Layout struct {
		Await token.Token
	}
	Value ast.Expr
}

func ParseAwaitExpr(p *parser.Parser) (node *AwaitExpr, err error) {
	node = &AwaitExpr{}
	if node.Layout.Await, err = p.Expect(AWAIT); err != nil {
		return
	}
	if node.Value, err = js.ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintAwaitExpr(pr *printer.Printer, node *AwaitExpr) error {
	pr.Print(node.Layout.Await)
	pr.Space().Print(node.Value)
	return nil
}
