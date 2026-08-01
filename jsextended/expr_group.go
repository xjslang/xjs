package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type GroupExpr struct {
	ast.BaseExpr
	Layout struct {
		Lparen token.Token
		Rparen token.Token
	}
	Value ast.Expr
}

func ParseGroupExpr(p *parser.Parser) (node *GroupExpr, err error) {
	node = &GroupExpr{}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if p.CurrentToken.Type != token.RPAREN {
		if node.Value, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	return node, nil
}

func PrintGroupExpr(pr *printer.Printer, node *GroupExpr) error {
	pr.Print(node.Layout.Lparen)
	if node.Value != nil {
		pr.Print(node.Value)
	}
	pr.Print(node.Layout.Rparen)
	return nil
}
