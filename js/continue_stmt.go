package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var CONTINUE = token.RegisterType("continue")

type ContinueStmt struct {
	ast.BaseStmt
	Layout struct {
		Continue token.Token
		Semi     token.Token
	}
	Label *Ident
}

func ParseContinueStmt(p *parser.Parser) (_ *ContinueStmt, err error) {
	node := &ContinueStmt{}
	if node.Layout.Continue, err = p.Expect(CONTINUE); err != nil {
		return
	}
	if p.CurrentToken.Type == token.IDENT && !p.CurrentToken.AfterNewline {
		if node.Label, err = ParseIdent(p); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintContinueStmt(p *printer.Printer, node *ContinueStmt) {
	p.LnPrint(node.Layout.Continue)
	if node.Label != nil {
		p.SpPrint(node.Label)
	}
	p.Print(node.Layout.Semi)
}
