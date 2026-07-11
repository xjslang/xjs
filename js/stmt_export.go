package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var EXPORT = token.RegisterType("export")

type ExportStmt struct {
	ast.BaseStmt
	Layout struct {
		Export token.Token
		Lbrace token.Token
		Rbrace token.Token
		Semi   token.Token
	}
	Decl    ast.Decl
	Exports []*ExportNode
}

type ExportNode struct {
	ast.BaseNode
	Layout struct {
		As token.Token
	}
	Name  *Ident
	Alias *Ident
}

func ParseExportStmt(p *parser.Parser) (node *ExportStmt, err error) {
	node = &ExportStmt{}
	if node.Layout.Export, err = p.Expect(EXPORT); err != nil {
		return
	}
	if p.CurrentToken.Type == token.LBRACE {
		node.Layout.Lbrace = p.CurrentToken
		p.AdvanceToken()
		for p.CurrentToken.Type != token.RBRACE {
			n := &ExportNode{}
			if n.Name, err = ParseIdent(p); err != nil {
				return
			}
			if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "as" {
				n.Layout.As = p.CurrentToken
				p.AdvanceToken()
				if n.Alias, err = ParseIdent(p); err != nil {
					return
				}
			}
			node.Exports = append(node.Exports, n)
			if p.CurrentToken.Type != token.COMMA {
				break
			}
			p.AdvanceToken()
		}
		if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
			return
		}
		if node.Layout.Semi, err = ExpectSemi(p); err != nil {
			return
		}
	} else {
		tok := p.CurrentToken
		var stmt ast.Stmt
		if stmt, err = p.ParseStmt(); err != nil {
			return
		}
		if v, ok := stmt.(ast.Decl); ok {
			node.Decl = v
		} else {
			err = p.ErrorAt(tok, "declaration expected")
		}
	}
	return
}

func PrintExportStmt(pr *printer.Printer, node *ExportStmt) error {
	pr.Line().Print(node.Layout.Export)
	if node.Decl != nil {
		pr.Space().Print(node.Decl)
	} else {
		pr.Space().Print(node.Layout.Lbrace)
		pr.IncreaseIndent()
		for i, item := range node.Exports {
			if i > 0 {
				pr.Print(',')
			}
			pr.Space().Print(item.Name)
			if item.Alias != nil {
				pr.Space().Print(item.Layout.As)
				pr.Space().Print(item.Alias)
			}
		}
		pr.DecreaseIndent()
		if len(node.Exports) > 0 {
			pr.Space()
		}
		pr.Print(node.Layout.Rbrace, node.Layout.Semi)
	}
	return nil
}
