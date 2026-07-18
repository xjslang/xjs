package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type SequenceExpr struct {
	ast.BaseExpr
	Layout struct {
		Lparen token.Token
		Rparen token.Token
	}
	Values []ast.Expr
}

func ParseSequenceExpr(p *parser.Parser) (node *SequenceExpr, err error) {
	node = &SequenceExpr{}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	for {
		if p.CurrentToken.Type == token.RPAREN {
			break
		}
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
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	return
}

func PrintSequenceExpr(pr *printer.Printer, node *SequenceExpr) error {
	pr.Print(node.Layout.Lparen)
	pr.IncreaseIndent()
	for i, val := range node.Values {
		if i > 0 {
			pr.Print(",")
			pr.Space()
		}
		pr.Print(val)
	}
	pr.DecreaseIndent()
	pr.Print(node.Layout.Rparen)
	return nil
}
