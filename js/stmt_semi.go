package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type SemiStmt struct {
	ast.BaseStmt
	Layout struct {
		Semi token.Token
	}
}

func (node *SemiStmt) SelfClosing() bool {
	return true
}

func ParseSemiStmt(p *parser.Parser) (node *SemiStmt, err error) {
	node = &SemiStmt{}
	if node.Layout.Semi, err = expectSemi(p); err != nil {
		return
	}
	return node, nil
}

func PrintSemiStmt(p *printer.Printer, node *SemiStmt) {
	p.Print(node.Layout.Semi)
}

func expectSemi(p *parser.Parser) (token.Token, error) {
	tok := p.CurrentToken
	if tok.Type == token.SEMICOLON {
		p.AdvanceToken()
		return tok, nil
	}
	if tok.AfterNewline {
		tok = token.Token{Type: token.SEMICOLON, Literal: token.SEMICOLON.String(), Position: tok.Position}
		return tok, nil
	}
	return tok, parser.NewErrorAtToken(p.CurrentToken, "; or newline expected")
}
