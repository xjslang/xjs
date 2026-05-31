package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ParenExpr struct {
	LparenToken token.Token
	RparenToken token.Token

	Value ast.Node
}

func (node *ParenExpr) Type() string {
	return "ParenExpr"
}

func ParseParenExpr(p *parser.Parser) (_ *ParenExpr, err error) {
	node := &ParenExpr{}
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Value, err = p.ParseExpr(); err != nil {
		return
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	return node, nil
}

func PrintParenExpr(p *printer.Printer, node *ParenExpr) {
	p.Print(node.LparenToken, node.Value, node.RparenToken)
}
