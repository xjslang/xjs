package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var LET = token.RegisterType("let")

type LetStmt struct {
	LetToken    token.Token
	AssignToken token.Token
	SemiToken   token.Token

	Name  token.Token
	Value ast.Node
}

func (node *LetStmt) Type() string {
	return "LetStmt"
}

func ParseLetStmt(p *parser.Parser) (_ *LetStmt, err error) {
	node := &LetStmt{}
	if node.LetToken, err = p.Expect(LET); err != nil {
		return
	}
	if node.Name, err = p.Expect(token.IDENT); err != nil {
		return
	}
	if node.AssignToken, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintLetStmt(p *printer.Printer, node *LetStmt) {
	p.LnPrint(node.LetToken)
	p.SpPrint(node.Name, node.AssignToken, node.Value)
	p.Print(node.SemiToken)
}

func LetStmtPlugin(b *builder.Builder) {
	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "let" {
			tok.Type = LET
		}
		return tok
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == LET {
			return ParseLetStmt(p)
		}
		return next()
	})
}
