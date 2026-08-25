package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var EXPORT = token.RegisterType("EXPORT", "export")

type ExportStmt struct {
	ast.BaseStmt
	Layout struct {
		Export token.Token
	}
	Variant ast.Stmt
}

type WildcharExportStmt struct {
	ast.BaseStmt
	Layout struct {
		Multiply token.Token
		As       token.Token
		From     token.Token
		Semi     token.Token
	}
	Alias *Ident
	From  token.Token
}

type DefaultExportStmt struct {
	ast.BaseStmt
	Layout struct {
		Default token.Token
		Semi    token.Token
	}
	Stmt ast.Stmt
}

type NamedListExportStmt struct {
	ast.BaseStmt
	Layout struct {
		Lbrace token.Token
		Rbrace token.Token
		From   token.Token
		Semi   token.Token
	}
	IsDefault  bool
	isReExport bool
	Items      []*NamedExportItem
	From       token.Token
}

type NamedExportStmt struct {
	ast.BaseStmt
	Stmt ast.Stmt
}

type NamedExportItem struct {
	ast.BaseNode
	Layout struct {
		As token.Token
	}
	Name, Alias *Ident
}

func ParseExportStmt(p *parser.Parser) (node *ExportStmt, err error) {
	node = &ExportStmt{}
	if node.Layout.Export, err = p.Expect(EXPORT); err != nil {
		return
	}
	if node.Variant, err = parser.Switch(p, func(p *parser.Parser) (ast.Stmt, error) {
		return parseWildcharExport(p)
	}, func(p *parser.Parser) (ast.Stmt, error) {
		return parseDefaultExport(p)
	}, func(p *parser.Parser) (ast.Stmt, error) {
		return parseNamedListExport(p)
	}, func(p *parser.Parser) (ast.Stmt, error) {
		return parseNamedExport(p)
	}); err != nil {
		return
	}
	return
}

func parseWildcharExport(p *parser.Parser) (node *WildcharExportStmt, err error) {
	node = &WildcharExportStmt{}
	if node.Layout.Multiply, err = p.Expect(token.MULTIPLY); err != nil {
		return
	}
	if p.CurrentToken.Literal == "as" {
		node.Layout.As = p.CurrentToken
		p.AdvanceToken()
		if node.Alias, err = ParseIdent(p); err != nil {
			return
		}
	}
	if node.Layout.From, err = p.ExpectLiteral("from"); err != nil {
		return
	}
	if node.From, err = p.Expect(token.STRING); err != nil {
		return
	}
	if node.Layout.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return
}

func parseDefaultExport(p *parser.Parser) (node *DefaultExportStmt, err error) {
	node = &DefaultExportStmt{}
	if node.Layout.Default, err = p.Expect(DEFAULT); err != nil {
		return
	}
	if node.Stmt, err = ParseExprStmt(p); err != nil {
		return
	}
	return
}

func parseNamedListExport(p *parser.Parser) (node *NamedListExportStmt, err error) {
	node = &NamedListExportStmt{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for {
		if p.CurrentToken.Type == token.RBRACE || p.CurrentToken.Type == token.EOF {
			break
		}
		item := &NamedExportItem{}
		if p.CurrentToken.Type == DEFAULT {
			item.Name = &Ident{Token: p.CurrentToken}
			node.IsDefault = true
			p.AdvanceToken()
		} else if item.Name, err = ParseIdent(p); err != nil {
			return
		}
		if p.CurrentToken.Literal == "as" {
			item.Layout.As = p.CurrentToken
			p.AdvanceToken()
			if p.CurrentToken.Type == DEFAULT {
				item.Alias = &Ident{Token: p.CurrentToken}
				p.AdvanceToken()
			} else if item.Alias, err = ParseIdent(p); err != nil {
				return
			}
		}
		node.Items = append(node.Items, item)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
		p.AdvanceToken()
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	if p.CurrentToken.Literal == "from" {
		node.Layout.From = p.CurrentToken
		node.isReExport = true
		p.AdvanceToken()
		if node.From, err = p.Expect(token.STRING); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return
}

func parseNamedExport(p *parser.Parser) (node *NamedExportStmt, err error) {
	node = &NamedExportStmt{}
	tok := p.CurrentToken
	var stmt ast.Stmt
	if stmt, err = p.ParseStmt(); err != nil {
		return
	}
	if _, ok := stmt.(ast.Decl); !ok {
		err = p.ErrorAt(tok, "declaration expected")
		return
	}
	node.Stmt = stmt
	return
}

func PrintExportStmt(pr *printer.Printer, node *ExportStmt) error {
	pr.Line().Print(node.Layout.Export)
	switch v := node.Variant.(type) {
	case *NamedExportStmt:
		printNamedExportStmt(pr, v)
	case *NamedListExportStmt:
		printNamedListExportStmt(pr, v)
	case *DefaultExportStmt:
		printDefaultExportStmt(pr, v)
	case *WildcharExportStmt:
		printWildcharExportStmt(pr, v)
	}
	return nil
}

func printNamedExportStmt(pr *printer.Printer, node *NamedExportStmt) {
	pr.Space().Print(node.Stmt)
}

func printNamedListExportStmt(pr *printer.Printer, node *NamedListExportStmt) {
	pr.Space().Print(node.Layout.Lbrace)
	if len(node.Items) > 0 {
		pr.IncreaseIndent()
		for i, item := range node.Items {
			if i > 0 {
				pr.Print(',')
				pr.Space()
			}
			pr.Space().Print(item.Name)
			if item.Alias != nil {
				pr.Space().Print(item.Layout.As)
				pr.Space().Print(item.Alias)
			}
		}
		pr.DecreaseIndent()
		pr.Space()
	}
	pr.Print(node.Layout.Rbrace)
	if node.isReExport {
		pr.Space().Print(node.Layout.From)
		pr.Space().Print(node.From)
	}
	pr.Print(node.Layout.Semi)
}

func printDefaultExportStmt(pr *printer.Printer, node *DefaultExportStmt) {
	pr.Space().Print(node.Layout.Default)
	pr.Space().Print(node.Stmt)
}

func printWildcharExportStmt(pr *printer.Printer, node *WildcharExportStmt) {
	pr.Space().Print(node.Layout.Multiply)
	if node.Alias != nil {
		pr.Space().Print(node.Layout.As)
		pr.Space().Print(node.Alias)
	}
	pr.Space().Print(node.Layout.From)
	pr.Space().Print(node.From)
	pr.Print(node.Layout.Semi)
}
