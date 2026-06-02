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
	Tokens struct {
		Function token.Token
		Lparen   token.Token
		Rparen   token.Token
	}
	Name *Ident
	Body *BlockStmt
}

func ParseFunctionDecl(p *parser.Parser) (_ *FunctionDecl, err error) {
	node := &FunctionDecl{}
	if node.Tokens.Function, err = p.Expect(FUNCTION); err != nil {
		return
	}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Tokens.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Tokens.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Body, err = ParseBlockStmt(p); err != nil {
		return
	}
	return node, nil
}

func PrintFunctionDecl(p *printer.Printer, node *FunctionDecl) {
	p.LnPrint(node.Tokens.Function)
	p.SpPrint(node.Name)
	p.Print(node.Tokens.Lparen, node.Tokens.Rparen)
	p.SpPrint(node.Body)
}
