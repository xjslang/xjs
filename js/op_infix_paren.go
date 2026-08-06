package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type InfixParenOp struct {
	ast.BaseExpr
	Layout struct {
		Lparen token.Token
		Rparen token.Token
	}
	Callee ast.Expr
	Args   []ast.Expr
}

func ParseInfixParenOp(p *parser.Parser, left ast.Expr) (node *InfixParenOp, err error) {
	node = &InfixParenOp{Callee: left}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RPAREN {
		var val ast.Expr
		if val, err = ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
			return
		}
		node.Args = append(node.Args, val)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
		p.AdvanceToken()
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return nil, err
	}
	return node, nil
}

func PrintInfixParenOp(pr *printer.Printer, node *InfixParenOp) error {
	pr.Print(node.Callee, node.Layout.Lparen)
	for i, arg := range node.Args {
		if i > 0 {
			pr.Print(",")
			pr.Space()
		}
		pr.Print(arg)
	}
	pr.Print(node.Layout.Rparen)
	return nil
}
