package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var YIELD = token.RegisterType("yield")

type YieldStmt struct {
	ast.BaseStmt
	Layout struct {
		Yield    token.Token
		Multiply token.Token
		Semi     token.Token
	}
	IsDelegating bool
	Expr         ast.Expr
}

func ParseYieldStmt(p *parser.Parser) (node *YieldStmt, err error) {
	node = &YieldStmt{}
	if node.Layout.Yield, err = p.Expect(YIELD); err != nil {
		return
	}
	if node.IsDelegating = p.CurrentToken.Type == token.MULTIPLY; node.IsDelegating {
		node.Layout.Multiply = p.CurrentToken
		p.AdvanceToken()
	}
	var semi token.Token
	if semi, err = js.ExpectSemi(p); err == nil {
		if node.IsDelegating {
			err = p.Error("expression expected")
			return
		}
		node.Layout.Semi = semi
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

func PrintYieldStmt(pr *printer.Printer, node *YieldStmt) error {
	pr.Line().Print(node.Layout.Yield)
	if node.IsDelegating {
		pr.Print(node.Layout.Multiply)
	}
	if node.Expr != nil {
		pr.Space().Print(node.Expr)
	}
	pr.Print(node.Layout.Semi)
	return nil
}
