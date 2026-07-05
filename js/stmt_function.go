package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FUNCTION = token.RegisterType("function")

type FunctionDecl struct {
	ast.BaseStmt
	Layout struct {
		Function token.Token
		Lparen   token.Token
		Rparen   token.Token
	}
	Name   *Ident
	Params []*Ident
	Body   *BlockStmt
}

func (node *FunctionDecl) SelfClosing() bool {
	return true
}

func ParseFunctionDecl(p *parser.Parser) (node *FunctionDecl, err error) {
	node = &FunctionDecl{}
	if node.Layout.Function, err = p.Expect(FUNCTION); err != nil {
		return
	}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RPAREN {
		var name *Ident
		if name, err = ParseIdent(p); err != nil {
			return
		}
		node.Params = append(node.Params, name)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
		p.AdvanceToken()
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Body, err = ParseBlockStmt(p); err != nil {
		return
	}
	return node, nil
}

func PrintFunctionDecl(p *printer.Printer, node *FunctionDecl) {
	p.Line().Print(node.Layout.Function)
	p.Space().Print(node.Name)
	p.Print(node.Layout.Lparen)
	p.IncreaseIndent()
	for i, param := range node.Params {
		if i > 0 {
			p.Print(",")
			p.Space()
		}
		p.Print(param)
	}
	p.DecreaseIndent()
	p.Print(node.Layout.Rparen)
	p.Space().Print(node.Body)
}
