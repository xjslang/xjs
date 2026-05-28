package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FUNCTION = token.RegisterType("function")

type Function struct {
	FunctionToken token.Token
	LparenToken   token.Token
	RparenToken   token.Token

	Name token.Token
	Body *ast.Block
}

func (node *Function) Type() string {
	return "Function"
}

func ParseFunction(p *parser.Parser) (_ *Function, err error) {
	node := &Function{}
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
	if node.Body, err = parser.ParseBlock(p); err != nil {
		return
	}
	return node, nil
}

func PrintFunction(p *printer.Printer, node *Function) {
	p.LnPrint(node.FunctionToken)
	p.SpPrint(node.Name)
	p.Print(node.LparenToken, node.RparenToken)
	p.SpPrint(node.Body)
}
