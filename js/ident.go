package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
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

func PrintIdent(p *printer.Printer, node *Ident) {
	p.Print(node.Token)
}
