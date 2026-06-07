package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type FunctionExpr struct {
	ast.BaseExpr
	Layout struct {
		Function token.Token
		Lparen   token.Token
		Rparen   token.Token
	}
	Params []*Ident
	Body   *BlockStmt
}

func ParseFunctionExpr(p *parser.Parser) (_ *FunctionExpr, err error) {
	node := &FunctionExpr{}
	if node.Layout.Function, err = p.Expect(FUNCTION); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if p.CurrentToken.Type != token.RPAREN {
		for {
			var name *Ident
			if name, err = ParseIdent(p); err != nil {
				return
			}
			node.Params = append(node.Params, name)
			if p.CurrentToken.Type == token.RPAREN {
				break
			}
			if _, err = p.Expect(token.COMMA); err != nil {
				return
			}
		}
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Body, err = ParseBlockStmt(p); err != nil {
		return
	}
	return node, nil
}

func PrintFunctionExpr(p *printer.Printer, node *FunctionExpr) {
	p.Print(node.Layout.Function)
	p.SpPrint(node.Layout.Lparen)
	p.IncreaseIndent()
	for i, param := range node.Params {
		if i > 0 {
			p.Print(",")
			p.EnsureSpace()
		}
		p.Print(param)
	}
	p.DecreaseIndent()
	p.Print(node.Layout.Rparen)
	p.SpPrint(node.Body)
}
