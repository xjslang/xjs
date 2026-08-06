package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type InfixCommaOp struct {
	ast.BaseExpr
	Values []ast.Expr
}

func ParseInfixCommaOp(p *parser.Parser, left ast.Expr) (node *InfixCommaOp, err error) {
	node = &InfixCommaOp{}
	node.Values = append(node.Values, left)
	if _, err = p.Expect(token.COMMA); err != nil {
		return
	}
	for {
		var val ast.Expr
		if val, err = ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
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

func PrintInfixCommaOp(pr *printer.Printer, node *InfixCommaOp) error {
	for i, val := range node.Values {
		if i > 0 {
			pr.Print(",")
			pr.Space()
		}
		pr.Print(val)
	}
	return nil
}
