package js

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type MemberExpr struct {
	ast.BaseExpr
	Layout struct {
		Dot token.Token
	}
	Left  ast.Expr
	Right *Ident
}

func ParseMemberExpr(p *parser.Parser, left ast.Expr) (node *MemberExpr, err error) {
	node = &MemberExpr{Left: left}
	if node.Layout.Dot, err = p.Expect(token.DOT); err != nil {
		return
	}
	if node.Right, err = parseMemberKey(p); err != nil {
		return
	}
	return
}

func parseMemberKey(p *parser.Parser) (node *Ident, err error) {
	node = &Ident{}
	if node.Token, err = expectMemberKey(p); err != nil {
		return
	}
	return node, nil
}

func expectMemberKey(p *parser.Parser) (token.Token, error) {
	tok := p.CurrentToken
	if r, s := utf8.DecodeRuneInString(tok.Literal); s == 0 || !scanner.IsLetter(r) {
		return tok, p.Error("key expected")
	}
	p.AdvanceToken()
	return tok, nil
}

func PrintMemberExpr(p *printer.Printer, node *MemberExpr) {
	p.Print(node.Left, node.Layout.Dot, node.Right)
}
