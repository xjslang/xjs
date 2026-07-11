package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FUNCTION = token.RegisterType("function")

type FunctionDecl struct {
	ast.BaseDecl
	Layout struct {
		Function token.Token
		Lparen   token.Token
		Rparen   token.Token
	}
	Name   *Ident
	Params []*Ident
	Body   *BlockStmt
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

func PrintFunctionDecl(pr *printer.Printer, node *FunctionDecl) error {
	pr.Line().Print(node.Layout.Function)
	pr.Space().Print(node.Name)
	pr.Print(node.Layout.Lparen)
	pr.IncreaseIndent()
	for i, param := range node.Params {
		if i > 0 {
			pr.Print(",")
			pr.Space()
		}
		pr.Print(param)
	}
	pr.DecreaseIndent()
	pr.Print(node.Layout.Rparen)
	pr.Space().Print(node.Body)
	return nil
}
