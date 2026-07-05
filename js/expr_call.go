package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type CallExpr struct {
	ast.BaseExpr
	Layout struct {
		Lparen token.Token
		Rparen token.Token
	}
	Callee ast.Expr
	Args   []ast.Expr
}

func ParseCallExpr(p *parser.Parser, left ast.Expr) (node *CallExpr, err error) {
	node = &CallExpr{Callee: left}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RPAREN {
		var val ast.Expr
		if val, err = p.ParseExpr(); err != nil {
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

func PrintCallExpr(p *printer.Printer, node *CallExpr) {
	p.Print(node.Callee, node.Layout.Lparen)
	for i, arg := range node.Args {
		if i > 0 {
			p.Print(",")
			p.Space()
		}
		p.Print(arg)
	}
	p.Print(node.Layout.Rparen)
}
