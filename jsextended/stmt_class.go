package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	CLASS   = token.RegisterType("class")
	EXTENDS = token.RegisterType("extends")
)

type ClassProperty struct {
	ast.BaseDecl
	Layout struct {
		Assign token.Token
		Semi   token.Token
	}
	Default ast.Expr
}

type ClassMethod struct {
	ast.BaseDecl
	Params *FunctionParams
	Body   *js.BlockStmt
}

type StaticInitializer struct {
	ast.BaseNode
	Body *js.BlockStmt
}

type ClassMember struct {
	ast.BaseDecl
	Layout struct {
		Static token.Token
	}
	Static bool
	Flag   token.Token // get or set
	Name   ast.Node
	Decl   ast.Node
}

type ClassStmt struct {
	ast.BaseDecl
	Layout struct {
		Class   token.Token
		Extends token.Token
		Lbrace  token.Token
		Rbrace  token.Token
	}
	Name    *js.Ident
	Extends ast.Expr
	Members []*ClassMember
}

func ParseClassStmt(p *parser.Parser) (node *ClassStmt, err error) {
	node = &ClassStmt{}
	isMethodName := func(typ token.Type) bool {
		return typ == token.IDENT || typ == token.NUMBER || typ == token.STRING || typ == token.LBRACKET
	}
	if node.Layout.Class, err = p.Expect(CLASS); err != nil {
		return
	}
	if node.Name, err = ParseIdent(p); err != nil {
		return
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
				m.Static = isMethodName(p.PeekToken.Type)
			}
		}
		if m.Static {
			m.Layout.Static = p.CurrentToken
			p.AdvanceToken()
		}
		if p.CurrentToken.Type == token.IDENT && isMethodName(p.PeekToken.Type) {
			if lit := p.CurrentToken.Literal; lit == "get" || lit == "set" {
				m.Flag = p.CurrentToken
				p.AdvanceToken()
			}
		}
		if m.Name, err = parseMethodName(p); err != nil {
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

func parseMethodName(p *parser.Parser) (_ ast.Node, err error) {
	switch p.CurrentToken.Type {
	case token.STRING, token.NUMBER:
		tok := p.CurrentToken
		p.AdvanceToken()
		return &js.Literal{Value: tok}, nil
	case token.LBRACKET:
		return js.ParseComputedExpr(p)
	default:
		return ParseIdent(p)
	}
}

func parseProperty(p *parser.Parser) (node *ClassProperty, err error) {
	node = &ClassProperty{}
	if p.CurrentToken.Type == token.ASSIGN {
		node.Layout.Assign = p.CurrentToken
		p.AdvanceToken()
		if node.Default, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = js.ExpectSemi(p); err != nil {
		return
	}
	return
}

func parseMethod(p *parser.Parser) (node *ClassMethod, err error) {
	node = &ClassMethod{}
	if node.Params, err = ParseFunctionParams(p); err != nil {
		return
	}
	if node.Body, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	return
}

func parseStaticInitializer(p *parser.Parser) (node *StaticInitializer, err error) {
	node = &StaticInitializer{}
	if node.Body, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	return
}

func PrintClassStmt(pr *printer.Printer, node *ClassStmt) (err error) {
	pr.Line().Print(node.Layout.Class)
	pr.Space().Print(node.Name)
	if node.Extends != nil {
		pr.Space().Print(node.Layout.Extends)
		pr.Space().Print(node.Extends)
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
			switch n := m.Name.(type) {
			case *js.Ident, *js.Literal:
				pr.Print(m.Name)
			case *js.ComputedExpr:
				pr.Print(n.Layout.Lbracket, n.Expr, n.Layout.Rbracket)
			default:
				pr.Print('[', m.Name, ']')
			}
			if err = PrintFunctionParams(pr, v.Params); err != nil {
				return
			}
			pr.Space().Print(v.Body)
		case *ClassProperty:
			if len(m.Flag.Literal) > 0 {
				return pr.Error("get/set are reserved for methods")
			}
			switch n := m.Name.(type) {
			case *js.ComputedExpr:
				pr.Print(n.Layout.Lbracket, n.Expr, n.Layout.Rbracket)
			default:
				pr.Print(m.Name)
			}
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
