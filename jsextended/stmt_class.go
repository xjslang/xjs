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

type ClassMember struct {
	ast.BaseDecl
	Layout struct {
		Static token.Token
	}
	Static bool
	Name   *js.Ident
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
	Name        *js.Ident
	ParentClass *js.Ident
	Members     []*ClassMember
}

func ParseClassStmt(p *parser.Parser) (node *ClassStmt, err error) {
	node = &ClassStmt{}
	if node.Layout.Class, err = p.Expect(CLASS); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
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
			m.Static = p.PeekToken.Type == token.IDENT && p.PeekToken.Type != token.LPAREN
		}
		if m.Static {
			m.Layout.Static = p.CurrentToken
			p.AdvanceToken()
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

func PrintClassStmt(pr *printer.Printer, node *ClassStmt) (err error) {
	pr.Line().Print(node.Layout.Class)
	pr.Space().Print(node.Name)
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
		pr.Print(m.Name)
		switch v := m.Decl.(type) {
		case *ClassMethod:
			if err = PrintFunctionParams(pr, v.Params); err != nil {
				return
			}
			pr.Space().Print(v.Body)
		case *ClassProperty:
			if v.Default != nil {
				pr.Space().Print(v.Layout.Assign)
				pr.Space().Print(v.Default)
			}
			pr.Print(v.Layout.Semi)
		default:
			return pr.Error("class member expected")
		}
	}
	pr.DecreaseIndent()
	pr.Line().Print(node.Layout.Rbrace)
	return
}
