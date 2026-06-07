package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FOR = token.RegisterType("for")

type ForStmt struct {
	ast.BaseStmt
	Layout struct {
		For    token.Token
		Lparen token.Token
		Semi1  token.Token
		Semi2  token.Token
		Rparen token.Token
	}
	Init  ast.Stmt
	Cond  ast.Stmt
	After ast.Stmt
	Then  ast.Stmt
}

func (node *ForStmt) SelfClosing() bool {
	if v, ok := node.Then.(SelfClosingStmt); ok {
		return v.SelfClosing()
	}
	return false
}

func ParseForStmt(p *parser.Parser) (_ *ForStmt, err error) {
	node := &ForStmt{}
	// for
	if node.Layout.For, err = p.Expect(FOR); err != nil {
		return
	}
	// (init; cond; after)
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Init, err = parseForClause(p); err != nil {
		return
	}
	if node.Layout.Semi1, err = p.Expect(token.SEMICOLON); err != nil {
		return
	}
	if node.Cond, err = parseForClause(p); err != nil {
		return
	}
	if node.Layout.Semi2, err = p.Expect(token.SEMICOLON); err != nil {
		return
	}
	if node.After, err = parseForClause(p); err != nil {
		return
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	// then
	if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	return node, nil
}

func PrintForStmt(p *printer.Printer, node *ForStmt) {
	// for
	p.LnPrint(node.Layout.For)
	p.SpPrint(node.Layout.Lparen)
	// (init; condition; after)
	p.IncreaseIndent()
	p.BsPrint(node.Init)
	p.Print(node.Layout.Semi1)
	p.SpPrint(node.Cond)
	p.Print(node.Layout.Semi2)
	p.SpPrint(node.After)
	p.DecreaseIndent()
	// then
	p.Print(node.Layout.Rparen)
	p.SpPrint(node.Then)
}

func parseForClause(p *parser.Parser) (node ast.Stmt, err error) {
	switch p.CurrentToken.Type {
	case LET:
		return ParseLetStmt(p)
	case token.IDENT:
		switch p.PeekToken.Type {
		case token.ASSIGN:
			return ParseAssignStmt(p)
		case token.INCREMENT:
			if !p.PeekToken.AfterNewline {
				return ParseIncStmt(p)
			}
		case token.DECREMENT:
			if !p.PeekToken.AfterNewline {
				return ParseDecStmt(p)
			}
		}
	}
	return ParseExprStmt(p)
}
