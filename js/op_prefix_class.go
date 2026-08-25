package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type PrefixClassOp struct {
	ast.BaseExpr
	Layout struct {
		Class   token.Token
		Extends token.Token
		Lbrace  token.Token
		Rbrace  token.Token
	}
	Name    *Ident
	Extends ast.Expr
	Members []ast.Stmt
}

func ParsePrefixClassOp(p *parser.Parser) (node *PrefixClassOp, err error) {
	node = &PrefixClassOp{}
	if node.Layout.Class, err = p.Expect(CLASS); err != nil {
		return
	}
	if p.CurrentToken.Type == token.IDENT {
		node.Name = &Ident{Token: p.CurrentToken}
		p.AdvanceToken()
	}
	if p.CurrentToken.Type == EXTENDS {
		node.Layout.Extends = p.CurrentToken
		p.AdvanceToken()
		if node.Extends, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for {
		for p.CurrentToken.Type == token.SEMICOLON {
			p.AdvanceToken()
		}
		if p.CurrentToken.Type == token.RBRACE || p.CurrentToken.Type == token.EOF {
			break
		}
		var member ast.Stmt
		if member, err = parser.Switch(p,
			func(p *parser.Parser) (ast.Stmt, error) { return parseClassInitializer(p) },
			func(p *parser.Parser) (ast.Stmt, error) { return parseClassField(p) },
			func(p *parser.Parser) (ast.Stmt, error) { return parseClassMethod(p) },
		); err != nil {
			return
		}
		node.Members = append(node.Members, member)
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	return
}

func PrintPrefixClassOp(pr *printer.Printer, node *PrefixClassOp) (err error) {
	pr.Print(node.Layout.Class)
	if node.Name != nil {
		pr.Space().Print(node.Name)
	}
	if node.Extends != nil {
		pr.Space().Print(node.Layout.Extends)
		pr.Space().Print(node.Extends)
	}
	pr.Space().Print(node.Layout.Lbrace)
	if len(node.Members) > 0 {
		pr.IncreaseIndent()
		for _, entry := range node.Members {
			switch v := entry.(type) {
			case *ClassField:
				err = printClassField(pr, v)
			case *ClassMethod:
				err = printClassMethod(pr, v)
			case *ClassInitializer:
				err = printClassInitializer(pr, v)
			default:
				err = pr.Error("class member expected")
			}
			if err != nil {
				return
			}
		}
		pr.DecreaseIndent()
	}
	pr.Line().Print(node.Layout.Rbrace)
	return
}
