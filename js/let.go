package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var LET = token.RegisterType("let")

type Let struct {
	LetToken    token.Token
	AssignToken token.Token
	SemiToken   token.Token

	Name  token.Token
	Value ast.Node
}

func (node *Let) Type() string {
	return "Let"
}

func ParseLet(p *parser.Parser) (_ *Let, err error) {
	node := &Let{}
	if node.LetToken, err = p.Expect(LET); err != nil {
		return
	}
	if node.Name, err = p.Expect(token.IDENT); err != nil {
		return
	}
	if node.AssignToken, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintLet(p *printer.Printer, node *Let) {
	p.LnPrint(node.LetToken)
	p.SpPrint(node.Name, node.AssignToken, node.Value)
	p.Print(node.SemiToken)
}
