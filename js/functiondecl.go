package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var FUNCTION = token.RegisterType("function")

type FunctionDecl struct {
	FunctionToken token.Token
	LparenToken   token.Token
	RparenToken   token.Token

	Name token.Token
	Body *ast.Block
}

func (node *FunctionDecl) Type() string {
	return "FunctionDecl"
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
	if node.Body, err = parser.ParseBlock(p); err != nil {
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

func FuncDeclPlugin(b *builder.Builder) {
	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "function" {
			tok.Type = FUNCTION
		}
		return tok
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == FUNCTION {
			return ParseFunctionDecl(p)
		}
		return next()
	})
}
