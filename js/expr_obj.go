package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ObjEntry struct {
	Key   *Ident
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
		// TODO: key can also be a number or a string
		if entry.Key, err = ParseMemberKey(p); err != nil {
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

func PrintObjExpr(p *printer.Printer, node *ObjExpr) {
	p.Print(node.Layout.Lbrace)
	if len(node.Entries) > 0 {
		p.IncreaseIndent()
		for i, entry := range node.Entries {
			if i > 0 {
				p.Print(",")
			}
			p.SpPrint(entry.Key)
			p.Print(":")
			p.SpPrint(entry.Value)
		}
		p.DecreaseIndent()
		p.EnsureSpace()
	}
	p.Print(node.Layout.Rbrace)
}
