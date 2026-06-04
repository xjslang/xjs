package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var BREAK = token.RegisterType("break")

type BreakStmt struct {
	ast.BaseStmt
	Layout struct {
		Break token.Token
		Semi  token.Token
	}
	Label *Ident
}

func ParseBreakStmt(p *parser.Parser) (_ *BreakStmt, err error) {
	node := &BreakStmt{}
	if node.Layout.Break, err = p.Expect(BREAK); err != nil {
		return
	}
	if p.CurrentToken.Type == token.IDENT {
		if node.Label, err = ParseIdent(p); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintBreakStmt(p *printer.Printer, node *BreakStmt) {
	p.LnPrint(node.Layout.Break)
	if node.Label != nil {
		p.SpPrint(node.Label)
	}
	p.Print(node.Layout.Semi)
}
