package js

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type ComputedExpr struct {
	ast.BaseExpr
	Layout struct {
		Lbracket token.Token
		Rbracket token.Token
	}
	Expr ast.Expr
}

type ObjEntry struct {
	Key   ast.Expr
	Value ast.Expr
}

type ObjExpr struct {
	ast.BaseExpr
	Layout struct {
		Lbrace token.Token
		Rbrace token.Token
	}
	Entries []ObjEntry
}

func ParseObjExpr(p *parser.Parser) (node *ObjExpr, err error) {
	node = &ObjExpr{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RBRACE {
		entry := ObjEntry{}
		if entry.Key, err = parseKeyExpr(p); err != nil {
			return
		}
		if _, err = p.Expect(token.COLON); err != nil {
			return
		}
		if entry.Value, err = p.ParseExpr(); err != nil {
			return
		}
		node.Entries = append(node.Entries, entry)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
		p.AdvanceToken()
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	return node, nil
}

func parseKeyExpr(p *parser.Parser) (node ast.Expr, err error) {
	switch p.CurrentToken.Type {
	case token.LBRACKET:
		n := &ComputedExpr{}
		n.Layout.Lbracket = p.CurrentToken
		p.AdvanceToken()
		if n.Expr, err = p.ParseExpr(); err != nil {
			return
		}
		if n.Layout.Rbracket, err = p.Expect(token.RBRACKET); err != nil {
			return
		}
		node = n
	case STRING, NUMBER:
		if node, err = ParseValue(p); err != nil {
			return
		}
	default:
		if r, s := utf8.DecodeRuneInString(p.CurrentToken.Literal); s > 0 && scanner.IsLetter(r) {
			p.CurrentToken.Type = token.IDENT
			node = &Variable{Name: p.CurrentToken}
			p.AdvanceToken()
		} else {
			err = p.Error("key expected")
			return
		}
	}
	return
}

func PrintObjExpr(p *printer.Printer, node *ObjExpr) {
	p.Print(node.Layout.Lbrace)
	if len(node.Entries) > 0 {
		p.IncreaseIndent()
		for i, entry := range node.Entries {
			if i > 0 {
				p.Print(",")
			}
			switch v := entry.Key.(type) {
			case *ComputedExpr:
				p.SpPrint(v.Layout.Lbracket).Print(v.Expr, v.Layout.Rbracket)
			default:
				p.SpPrint(v)
			}
			p.Print(":")
			p.SpPrint(entry.Value)
		}
		p.DecreaseIndent()
		p.EnsureSpace()
	}
	p.Print(node.Layout.Rbrace)
}
