package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FUNCTION = token.RegisterType("function")

type FunctionDecl struct {
	ast.StmtNode
	FunctionToken token.Token
	LparenToken   token.Token
	RparenToken   token.Token

	Name token.Token
	Body *BlockStmt
}

func ParseFunctionDecl(p *parser.Parser) (_ *FunctionDecl, err error) {
	node := &FunctionDecl{}
	if node.FunctionToken, err = p.Expect(FUNCTION); err != nil {
		return
	}
	if node.Name, err = p.Expect(token.IDENT); err != nil {
		return
	}
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Body, err = ParseBlockStmt(p); err != nil {
		return
	}
	return node, nil
}

func PrintFunctionDecl(p *printer.Printer, node *FunctionDecl) {
	p.LnPrint(node.FunctionToken)
	p.SpPrint(node.Name)
	p.Print(node.LparenToken, node.RparenToken)
	p.SpPrint(node.Body)
}
