package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ClassExpr struct {
	ast.BaseExpr
	Layout struct {
		Class   token.Token
		Extends token.Token
		Lbrace  token.Token
		Rbrace  token.Token
	}
	Name        *js.Ident
	ParentClass *js.Ident
	Members     []*ClassMember
}

func ParseClassExpr(p *parser.Parser) (node *ClassExpr, err error) {
	node = &ClassExpr{}
	if node.Layout.Class, err = p.Expect(CLASS); err != nil {
		return
	}
	if p.CurrentToken.Type == token.IDENT {
		node.Name = &js.Ident{Token: p.CurrentToken}
		p.AdvanceToken()
	}
	if p.CurrentToken.Type == EXTENDS {
		node.Layout.Extends = p.CurrentToken
		p.AdvanceToken()
		if node.ParentClass, err = js.ParseIdent(p); err != nil {
			return
		}
	}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RBRACE && p.CurrentToken.Type != token.EOF {
		m := &ClassMember{}
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "static" {
			if p.PeekToken.Type == token.LBRACE {
				m.Layout.Static = p.CurrentToken
				p.AdvanceToken()
				if m.Decl, err = parseStaticInitializer(p); err != nil {
					return
				}
				node.Members = append(node.Members, m)
				continue
			} else {
				m.Static = p.PeekToken.Type == token.IDENT
			}
		}
		if m.Static {
			m.Layout.Static = p.CurrentToken
			p.AdvanceToken()
		}
		if p.CurrentToken.Type == token.IDENT && p.PeekToken.Type == token.IDENT {
			if lit := p.CurrentToken.Literal; lit == "get" || lit == "set" {
				m.Flag = p.CurrentToken
				p.AdvanceToken()
			}
		}
		if m.Name, err = js.ParseIdent(p); err != nil {
			return
		}
		if p.CurrentToken.Type == token.LPAREN {
			m.Decl, err = parseMethod(p)
		} else {
			m.Decl, err = parseProperty(p)
		}
		if err != nil {
			return
		}
		node.Members = append(node.Members, m)
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	return
}

func PrintClassExpr(pr *printer.Printer, node *ClassExpr) (err error) {
	pr.Print(node.Layout.Class)
	if node.Name != nil {
		pr.Space().Print(node.Name)
	}
	if node.ParentClass != nil {
		pr.Space().Print(node.Layout.Extends)
		pr.Space().Print(node.ParentClass)
	}
	pr.Space().Print(node.Layout.Lbrace)
	pr.IncreaseIndent()
	for _, m := range node.Members {
		pr.Line()
		if m.Static {
			pr.Print(m.Layout.Static)
			pr.Space()
		}
		switch v := m.Decl.(type) {
		case *ClassMethod:
			if len(m.Flag.Literal) > 0 {
				pr.Print(m.Flag)
				pr.Space()
			}
			pr.Print(m.Name)
			if err = PrintFunctionParams(pr, v.Params); err != nil {
				return
			}
			pr.Space().Print(v.Body)
		case *ClassProperty:
			if len(m.Flag.Literal) > 0 {
				return pr.Error("get/set are reserved for methods")
			}
			pr.Print(m.Name)
			if v.Default != nil {
				pr.Space().Print(v.Layout.Assign)
				pr.Space().Print(v.Default)
			}
			pr.Print(v.Layout.Semi)
		case *StaticInitializer:
			pr.Print(m.Layout.Static)
			pr.Space().Print(v.Body)
		default:
			return pr.Error("class member expected")
		}
	}
	pr.DecreaseIndent()
	pr.Line().Print(node.Layout.Rbrace)
	return
}
