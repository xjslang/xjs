package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var YIELD = token.RegisterType("YIELD", "yield")

type PrefixYieldOp struct {
	ast.BaseExpr
	Layout struct {
		Yield    token.Token
		Multiply token.Token
	}
	IsDelegating bool
	Expr         ast.Expr
}

func ParsePrefixYieldOp(p *parser.Parser) (node *PrefixYieldOp, err error) {
	node = &PrefixYieldOp{}
	if node.Layout.Yield, err = p.Expect(YIELD); err != nil {
		return
	}
	if node.IsDelegating = p.CurrentToken.Type == token.MULTIPLY; node.IsDelegating {
		node.Layout.Multiply = p.CurrentToken
		p.AdvanceToken()
	}
	if node.IsDelegating || p.CurrentToken.Type != token.COMMA && !js.IsSemi(p.CurrentToken) {
		if node.Expr, err = js.ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
			return
		}
	}
	return
}

func PrintPrefixYieldOp(pr *printer.Printer, node *PrefixYieldOp) error {
	pr.Print(node.Layout.Yield)
	if node.IsDelegating {
		pr.Print(node.Layout.Multiply)
	}
	if node.Expr != nil {
		pr.Space().Print(node.Expr)
	}
	return nil
}
