package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ArrayExpr struct {
	ast.BaseExpr
	Layout struct {
		Lbracket token.Token
		Rbracket token.Token
	}
	Values []ast.Expr
}

func ParseArrayExpr(p *parser.Parser) (_ *ArrayExpr, err error) {
	node := &ArrayExpr{}
	if node.Layout.Lbracket, err = p.Expect(token.LBRACKET); err != nil {
		return
	}
	if p.CurrentToken.Type != token.RBRACKET {
		for {
			var val ast.Expr
			if val, err = p.ParseExpr(); err != nil {
				return
			}
			node.Values = append(node.Values, val)
			if p.CurrentToken.Type == token.RBRACKET {
				break
			}
			if _, err = p.Expect(token.COMMA); err != nil {
				return
			}
		}
	}
	if node.Layout.Rbracket, err = p.Expect(token.RBRACKET); err != nil {
		return
	}
	return node, nil
}

func PrintArrayExpr(p *printer.Printer, node *ArrayExpr) {
	p.Print(node.Layout.Lbracket)
	if len(node.Values) > 0 {
		p.IncreaseIndent()
		for i, val := range node.Values {
			if i > 0 {
				p.Print(",")
				p.EnsureSpace()
			}
			p.Print(val)
		}
		p.DecreaseIndent()
	}
	p.Print(node.Layout.Rbracket)
}
