package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	PLUS_ASSIGN     = token.RegisterType("+=")
	MINUS_ASSIGN    = token.RegisterType("-=")
	MULTIPLY_ASSIGN = token.RegisterType("*=")
	DIVIDE_ASSIGN   = token.RegisterType("/=")
	MODULO_ASSIGN   = token.RegisterType("%=")
)

type CompoundAssignExpr struct {
	ast.BaseExpr
	Layout struct {
		Op token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseCompoundAssignExpr(p *parser.Parser, left ast.Expr) (node *CompoundAssignExpr, err error) {
	node = &CompoundAssignExpr{Left: left}
	if typ := p.CurrentToken.Type; typ != PLUS_ASSIGN && typ != MINUS_ASSIGN && typ != MULTIPLY_ASSIGN && typ != DIVIDE_ASSIGN && typ != MODULO_ASSIGN {
		err = p.Error("compound assignment expected")
		return
	}
	node.Layout.Op = p.CurrentToken
	p.AdvanceToken()
	if node.Right, err = p.ParseExpr(); err != nil {
		return
	}
	return
}

func PrintCompoundAssignExpr(pr *printer.Printer, node *CompoundAssignExpr) error {
	pr.Print(node.Left)
	pr.Space().Print(node.Layout.Op)
	pr.Space().Print(node.Right)
	return nil
}
