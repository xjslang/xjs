package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var YIELD = token.RegisterType("yield")

type YieldExpr struct {
	ast.BaseExpr
	Layout struct {
		Yield    token.Token
		Multiply token.Token
	}
	IsDelegating bool
	Expr         ast.Expr
}

func ParseYieldExpr(p *parser.Parser) (node *YieldExpr, err error) {
	node = &YieldExpr{}
	if node.Layout.Yield, err = p.Expect(YIELD); err != nil {
		return
	}
	if node.IsDelegating = p.CurrentToken.Type == token.MULTIPLY; node.IsDelegating {
		node.Layout.Multiply = p.CurrentToken
		p.AdvanceToken()
	}
	if node.Expr, err = js.ParseRightExpr(p, token.LPAREN.Precedence()-1); err != nil {
		return
	}
	return
}

func PrintYieldExpr(pr *printer.Printer, node *YieldExpr) error {
	pr.Print(node.Layout.Yield)
	if node.IsDelegating {
		pr.Print(node.Layout.Multiply)
	}
	pr.Space().Print(node.Expr)
	return nil
}
