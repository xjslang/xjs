package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type EmptyStmt struct {
	ast.BaseStmt
	Layout struct {
		Semi token.Token
	}
}

func ParseEmptyStmt(p *parser.Parser) (_ *EmptyStmt, err error) {
	node := &EmptyStmt{}
	if node.Layout.Semi, err = p.Expect(token.SEMICOLON); err != nil {
		return
	}
	return node, nil
}

func PrintEmptyStmt(p *printer.Printer, node *EmptyStmt) {
	if len(node.Layout.Semi.LeadingTrivia) > 0 {
		p.PrintTrivia(node.Layout.Semi.LeadingTrivia)
	}
}
