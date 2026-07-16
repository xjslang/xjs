package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var IMPORT = token.RegisterType("import")

type ImportStmt struct {
	ast.BaseStmt
	Layout struct {
		Import   token.Token
		Lbrace   token.Token
		Rbrace   token.Token
		Multiply token.Token
		As       token.Token
		From     token.Token
		Semi     token.Token
	}
	Imports   []*ImportNode
	Namespace *Ident
	Default   *Ident
	Path      token.Token
}

type ImportNode struct {
	ast.BaseNode
	Layout struct {
		As token.Token
	}
	Name  *Ident
	Alias *Ident
}

func ParseImportStmt(p *parser.Parser) (node *ImportStmt, err error) {
	node = &ImportStmt{}
	if node.Layout.Import, err = p.Expect(IMPORT); err != nil {
		return
	}
	if p.CurrentToken.Type == token.STRING {
		// side-effects import:
		// import 'lib.js'
		node.Path = p.CurrentToken
		p.AdvanceToken()
	} else {
		switch p.CurrentToken.Type {
		case token.MULTIPLY:
			// namespace import:
			// import * as lib from 'lib.js'
			node.Layout.Multiply = p.CurrentToken
			p.AdvanceToken()
			if node.Layout.As, err = p.ExpectString("as"); err != nil {
				return
			}
			if node.Namespace, err = ParseIdent(p); err != nil {
				return
			}
		case token.IDENT:
			// default import:
			// import lib from 'lib.js'
			if node.Default, err = ParseIdent(p); err != nil {
				return
			}
		case token.LBRACE:
			// named imports:
			// import { c1, c2 as c3, c4 } from 'library'
			node.Layout.Lbrace = p.CurrentToken
			p.AdvanceToken()
			for p.CurrentToken.Type != token.RBRACE {
				e := &ImportNode{}
				if e.Name, err = ParseIdent(p); err != nil {
					return
				}
				if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "as" {
					e.Layout.As = p.CurrentToken
					p.AdvanceToken()
					if e.Alias, err = ParseIdent(p); err != nil {
						return
					}
				}
				node.Imports = append(node.Imports, e)
				if p.CurrentToken.Type != token.COMMA {
					break
				}
				p.AdvanceToken()
			}
			if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
				return
			}
		default:
			err = p.Error("syntax error")
			return
		}
		if node.Layout.From, err = p.ExpectString("from"); err != nil {
			return
		}
		if node.Path, err = p.Expect(token.STRING); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return
}

func PrintImportStmt(pr *printer.Printer, node *ImportStmt) error {
	pr.Line().Print(node.Layout.Import)
	if node.Namespace != nil {
		// namespace import
		pr.Space().Print(node.Layout.Multiply)
		pr.Space().Print(node.Layout.As)
		pr.Space().Print(node.Namespace)
		pr.Space().Print(node.Layout.From)
	} else if node.Default != nil {
		// default import
		pr.Space().Print(node.Default)
		pr.Space().Print(node.Layout.From)
	} else {
		// named exports
		pr.Space().Print(node.Layout.Lbrace)
		pr.IncreaseIndent()
		for i, export := range node.Imports {
			if i > 0 {
				pr.Print(',')
			}
			pr.Space().Print(export.Name)
			if export.Alias != nil {
				pr.Space().Print(export.Layout.As)
				pr.Space().Print(export.Alias)
			}
		}
		pr.DecreaseIndent()
		if len(node.Imports) > 0 {
			pr.Space()
		}
		pr.Print(node.Layout.Rbrace)
		pr.Space().Print(node.Layout.From)
	}
	pr.Space().Print(node.Path, node.Layout.Semi)
	return nil
}
