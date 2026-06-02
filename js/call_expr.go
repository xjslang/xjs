package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type CallExpr struct {
	ast.ExprNode
	LparenToken token.Token
	RparenToken token.Token

	Function  ast.Expr
	Arguments []ast.Expr
}

func (node *CallExpr) Type() string {
	return "CallExpr"
}

func ParseCallExpr(p *parser.Parser, leftVal ast.Expr) (_ *CallExpr, err error) {
	node := &CallExpr{Function: leftVal}
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if p.CurrentToken.Type != token.RPAREN {
		for {
			var val ast.Expr
			if val, err = p.ParseExpr(); err != nil {
				return
			}
			node.Arguments = append(node.Arguments, val)
			if p.CurrentToken.Type == token.RPAREN {
				break
			}
			if _, err = p.Expect(token.COMMA); err != nil {
				return
			}
		}
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return nil, err
	}
	return node, nil
}

func PrintCallExpr(p *printer.Printer, node *CallExpr) {
	p.Print(node.Function, node.LparenToken)
	for i, arg := range node.Arguments {
		if i > 0 {
			p.Print(",")
			p.EnsureSpace()
		}
		p.Print(arg)
	}
	p.Print(node.RparenToken)
}
