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

func ParseContinueStmt(p *parser.Parser) (node *ContinueStmt, err error) {
	node = &ContinueStmt{}
	if node.Layout.Continue, err = p.Expect(CONTINUE); err != nil {
		return
	}
	if p.CurrentToken.Type == token.IDENT && !p.CurrentToken.AfterNewline {
		if node.Label, err = ParseIdent(p); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return
}

func PrintContinueStmt(p *printer.Printer, node *ContinueStmt) error {
	p.Line().Print(node.Layout.Continue)
	if node.Label != nil {
		p.Space().Print(node.Label)
	}
	p.Print(node.Layout.Semi)
	return nil
}
