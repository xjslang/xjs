package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type SequenceExpr struct {
	ast.BaseExpr
	Values []ast.Expr
}

func ParseSequenceExpr(p *parser.Parser, left ast.Expr) (node *SequenceExpr, err error) {
	node = &SequenceExpr{}
	node.Values = append(node.Values, left)
	if _, err = p.Expect(token.COMMA); err != nil {
		return
	}
	for {
		var val ast.Expr
		if val, err = js.ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
			return
		}
		node.Values = append(node.Values, val)
		if p.CurrentToken.Type == token.COMMA {
			p.AdvanceToken()
			continue
		}
		break
	}
	return
}

func PrintSequenceExpr(pr *printer.Printer, node *SequenceExpr) error {
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
	return nil
}
