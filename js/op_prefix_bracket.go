package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type PrefixBracketOp struct {
	ast.BaseExpr
	Layout struct {
		Lbracket token.Token
		Rbracket token.Token
	}
	Values []ast.Expr
}

func ParsePrefixBracketOp(p *parser.Parser) (node *PrefixBracketOp, err error) {
	node = &PrefixBracketOp{}
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
		if p.CurrentToken.Type == token.RBRACKET {
			break
		}
		var val ast.Expr
		if val, err = ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
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

func PrintPrefixBracketOp(pr *printer.Printer, node *PrefixBracketOp) error {
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
