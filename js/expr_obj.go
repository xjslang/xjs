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
	Key   ast.Node
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
		switch p.CurrentToken.Type {
		case token.LBRACKET:
			if entry.Key, err = ParseComputedExpr(p); err != nil {
				return
			}
		case token.STRING, token.NUMBER:
			if entry.Key, err = ParseValue(p); err != nil {
				return
			}
		default:
			if entry.Key, err = ParseObjKey(p); err != nil {
				return
			}
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

func ParseObjKey(p *parser.Parser) (node *Ident, err error) {
	tok := p.CurrentToken
	if r, s := utf8.DecodeRuneInString(tok.Literal); s == 0 || !scanner.IsLetter(r) {
		err = p.Error("key expected")
		return
	}
	node = &Ident{Token: tok}
	p.AdvanceToken()
	return
}

func ParseComputedExpr(p *parser.Parser) (node *ComputedExpr, err error) {
	node = &ComputedExpr{}
	node.Layout.Lbracket = p.CurrentToken
	p.AdvanceToken()
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rbracket, err = p.Expect(token.RBRACKET); err != nil {
		return
	}
	return
}

func PrintObjExpr(pr *printer.Printer, node *ObjExpr) error {
	pr.Print(node.Layout.Lbrace)
	if len(node.Entries) > 0 {
		pr.IncreaseIndent()
		for i, entry := range node.Entries {
			if i > 0 {
				pr.Print(",")
			}
			switch v := entry.Key.(type) {
			case *ComputedExpr:
				pr.Space().Print(v.Layout.Lbracket)
				pr.Print(v.Expr, v.Layout.Rbracket)
			default:
				pr.Space().Print(v)
			}
			pr.Print(":")
			pr.Space().Print(entry.Value)
		}
		pr.DecreaseIndent()
		pr.Space()
	}
	pr.Print(node.Layout.Rbrace)
	return nil
}
