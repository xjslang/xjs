package js

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type Ident struct {
	ast.BaseNode
	token.Token
}

func ParseIdent(p *parser.Parser) (node *Ident, err error) {
	node = &Ident{}
	if node.Token, err = p.Expect(token.IDENT); err != nil {
		return
	}
	return node, nil
}

func ParseMemberKey(p *parser.Parser) (node *Ident, err error) {
	node = &Ident{}
	if node.Token, err = ExpectMemberKey(p); err != nil {
		return
	}
	return node, nil
}

func ExpectMemberKey(p *parser.Parser) (token.Token, error) {
	tok := p.CurrentToken
	if r, s := utf8.DecodeRuneInString(tok.Literal); s == 0 || !scanner.IsLetter(r) {
		return tok, p.Error("key expected")
	}
	p.AdvanceToken()
	return tok, nil
}

func PrintIdent(p *printer.Printer, node *Ident) {
	p.Print(node.Token)
}
