package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var AWAIT = token.RegisterType("AWAIT", "await")

type PrefixAwaitOp struct {
	ast.BaseExpr
	Layout struct {
		Await token.Token
	}
	Value ast.Expr
}

func ParsePrefixAwaitOp(p *parser.Parser) (node *PrefixAwaitOp, err error) {
	node = &PrefixAwaitOp{}
	if node.Layout.Await, err = p.Expect(AWAIT); err != nil {
		return
	}
	if node.Value, err = js.ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintPrefixAwaitOp(pr *printer.Printer, node *PrefixAwaitOp) error {
	pr.Print(node.Layout.Await)
	pr.Space().Print(node.Value)
	return nil
}
