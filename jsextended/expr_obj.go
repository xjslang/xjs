package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ObjEntry struct {
	Key     ast.Node
	Value   ast.Expr
	Default ast.Expr
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
			if entry.Key, err = js.ParseComputedExpr(p); err != nil {
				return
			}
		case token.STRING, token.NUMBER, SPREAD:
			if entry.Key, err = js.ParseValue(p); err != nil {
				return
			}
		default:
			if entry.Key, err = js.ParseObjKey(p); err != nil {
				return
			}
		}
		if p.CurrentToken.Type == token.COLON {
			p.AdvanceToken()
			if entry.Value, err = js.ParseRightExpr(p, token.ASSIGN.Precedence()); err != nil {
				return
			}
		}
		if p.CurrentToken.Type == token.ASSIGN {
			p.AdvanceToken()
			if entry.Default, err = p.ParseExpr(); err != nil {
				return
			}
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

func PrintObjExpr(pr *printer.Printer, node *ObjExpr) error {
	pr.Print(node.Layout.Lbrace)
	if len(node.Entries) > 0 {
		pr.IncreaseIndent()
		for i, entry := range node.Entries {
			if i > 0 {
				pr.Print(",")
			}
			switch v := entry.Key.(type) {
			case *js.ComputedExpr:
				pr.Space().Print(v.Layout.Lbracket, v.Expr, v.Layout.Rbracket)
			default:
				pr.Space().Print(v)
			}
			if entry.Value != nil {
				pr.Print(":")
				pr.Space().Print(entry.Value)
			}
			if entry.Default != nil {
				pr.Space().Print("=")
				pr.Space().Print(entry.Default)
			}
		}
		pr.DecreaseIndent()
		pr.Space()
	}
	pr.Print(node.Layout.Rbrace)
	return nil
}
