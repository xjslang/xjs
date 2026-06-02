package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ParenExpr struct {
	ast.ExprNode
	Tokens struct {
		Lparen token.Token
		Rparen token.Token
	}
	Value ast.Expr
}

func ParseParenExpr(p *parser.Parser) (_ *ParenExpr, err error) {
	node := &ParenExpr{}
	if node.Tokens.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Value, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Tokens.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	return node, nil
}

func PrintParenExpr(p *printer.Printer, node *ParenExpr) {
	p.Print(node.Tokens.Lparen, node.Value, node.Tokens.Rparen)
}
