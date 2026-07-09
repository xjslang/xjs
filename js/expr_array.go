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

func ParseArrayExpr(p *parser.Parser) (node *ArrayExpr, err error) {
	node = &ArrayExpr{}
	if node.Layout.Lbracket, err = p.Expect(token.LBRACKET); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RBRACKET {
		var val ast.Expr
		if val, err = p.ParseExpr(); err != nil {
			return
		}
		node.Values = append(node.Values, val)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
		p.AdvanceToken()
	}
	if node.Layout.Rbracket, err = p.Expect(token.RBRACKET); err != nil {
		return
	}
	return node, nil
}

func PrintArrayExpr(pr *printer.Printer, node *ArrayExpr) error {
	pr.Print(node.Layout.Lbracket)
	if len(node.Values) > 0 {
		pr.IncreaseIndent()
		for i, val := range node.Values {
			if i > 0 {
				pr.Print(",")
				pr.Space()
			}
			pr.Print(val)
		}
		pr.DecreaseIndent()
	}
	pr.Print(node.Layout.Rbracket)
	return nil
}
