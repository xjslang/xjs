package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

func ParseArrayExpr(p *parser.Parser) (node *js.ArrayExpr, err error) {
	node = &js.ArrayExpr{}
	if node.Layout.Lbracket, err = p.Expect(token.LBRACKET); err != nil {
		return
	}
	var prevElem ast.Node
	for p.CurrentToken.Type != token.RBRACKET {
		for {
			if prevElem != nil {
				if _, err = p.Expect(token.COMMA); err != nil {
					return
				}
			} else if p.CurrentToken.Type == token.COMMA {
				p.AdvanceToken()
			} else {
				break
			}
			if prevElem == nil {
				node.Values = append(node.Values, nil)
			}
			prevElem = nil
		}
		var val ast.Expr
		if val, err = p.ParseExpr(); err != nil {
			return
		}
		node.Values = append(node.Values, val)
		prevElem = val
	}
	if node.Layout.Rbracket, err = p.Expect(token.RBRACKET); err != nil {
		return
	}
	return node, nil
}

func PrintArrayExpr(pr *printer.Printer, node *js.ArrayExpr) error {
	pr.Print(node.Layout.Lbracket)
	if len(node.Values) > 0 {
		pr.IncreaseIndent()
		for i, val := range node.Values {
			if i > 0 {
				pr.Print(",")
				pr.Space()
			}
			if val != nil {
				pr.Print(val)
			} else {
				pr.Space()
			}
		}
		pr.DecreaseIndent()
	}
	pr.Print(node.Layout.Rbracket)
	return nil
}
