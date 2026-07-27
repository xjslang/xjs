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

type ClassMethodStmt struct {
	ast.BaseDecl
	Name   *js.Ident
	Params *FunctionParams
	Body   *js.BlockStmt
}

type ClassStmt struct {
	ast.BaseDecl
	Layout struct {
		Class   token.Token
		Extends token.Token
		Lbrace  token.Token
		Rbrace  token.Token
	}
	Name      *js.Ident
	BaseClass *js.Ident
	Methods   []*ClassMethodStmt
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
		if node.BaseClass, err = js.ParseIdent(p); err != nil {
			return
		}
	}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RBRACE && p.CurrentToken.Type != token.EOF {
		var m *ClassMethodStmt
		if m, err = parseMethod(p); err != nil {
			return
		}
		node.Methods = append(node.Methods, m)
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	return
}

func parseMethod(p *parser.Parser) (node *ClassMethodStmt, err error) {
	node = &ClassMethodStmt{}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
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
	if node.BaseClass != nil {
		pr.Space().Print(node.Layout.Extends)
		pr.Space().Print(node.BaseClass)
	}
	pr.Space().Print(node.Layout.Lbrace)
	pr.IncreaseIndent()
	for _, m := range node.Methods {
		pr.Line().Print(m.Name)
		if err = PrintFunctionParams(pr, m.Params); err != nil {
			return
		}
		pr.Space().Print(m.Body)
	}
	pr.DecreaseIndent()
	pr.Line().Print(node.Layout.Rbrace)
	return
}
