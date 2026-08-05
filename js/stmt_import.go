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
		Import token.Token
	}
	Variant ast.Stmt
}

type SideEffectImportStmt struct {
	ast.BaseStmt
	Layout struct {
		Semi token.Token
	}
	From token.Token
}

type DefaultImportStmt struct {
	ast.BaseStmt
	Layout struct {
		Comma token.Token
		From  token.Token
		Semi  token.Token
	}
	Name    *Ident
	Variant ast.Stmt
	From    token.Token
}

type WildcardImportStmt struct {
	ast.BaseStmt
	Layout struct {
		Multiply token.Token
		As       token.Token
		From     token.Token
		Semi     token.Token
	}
	Name *Ident
	From token.Token
}

type NamedListImportStmt struct {
	ast.BaseStmt
	Layout struct {
		Lbrace token.Token
		Rbrace token.Token
		From   token.Token
		Semi   token.Token
	}
	Items []*NamedImport
	From  token.Token
}

type NamedImport struct {
	ast.BaseNode
	Layout struct {
		As token.Token
	}
	Pattern ast.Node
	Alias   *Ident
}

func ParseImportStmt(p *parser.Parser) (node *ImportStmt, err error) {
	node = &ImportStmt{}
	if node.Layout.Import, err = p.Expect(IMPORT); err != nil {
		return
	}
	if node.Variant, err = parser.Switch(p, func(p *parser.Parser) (ast.Stmt, error) {
		return parseSideEffectImport(p)
	}, func(p *parser.Parser) (ast.Stmt, error) {
		return parseDefaultImport(p)
	}, func(p *parser.Parser) (ast.Stmt, error) {
		return parseWildcardImport(p)
	}, func(p *parser.Parser) (ast.Stmt, error) {
		return parseNamedListImport(p)
	}); err != nil {
		return
	}
	return
}

func parseSideEffectImport(p *parser.Parser) (node *SideEffectImportStmt, err error) {
	node = &SideEffectImportStmt{}
	if node.From, err = p.Expect(token.STRING); err != nil {
		return
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return
}

func parseWildcardImport(p *parser.Parser) (node *WildcardImportStmt, err error) {
	node = &WildcardImportStmt{}
	if node.Layout.Multiply, err = p.Expect(token.MULTIPLY); err != nil {
		return
	}
	if node.Layout.As, err = p.ExpectLiteral("as"); err != nil {
		return
	}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.From, err = p.ExpectLiteral("from"); err != nil {
		return
	}
	if node.From, err = p.Expect(token.STRING); err != nil {
		return
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return
}

func parseNamedListImport(p *parser.Parser) (node *NamedListImportStmt, err error) {
	node = &NamedListImportStmt{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for {
		if p.CurrentToken.Type == token.RBRACE || p.CurrentToken.Type == token.EOF {
			break
		}
		item := &NamedImport{}
		if p.CurrentToken.Type == DEFAULT {
			item.Pattern = &Ident{Token: p.CurrentToken}
			p.AdvanceToken()
		} else if item.Pattern, err = ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
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
	if node.Layout.From, err = p.ExpectLiteral("from"); err != nil {
		return
	}
	if node.From, err = p.Expect(token.STRING); err != nil {
		return
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return
}

func parseDefaultImport(p *parser.Parser) (node *DefaultImportStmt, err error) {
	node = &DefaultImportStmt{}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if p.CurrentToken.Type == token.COMMA {
		node.Layout.Comma = p.CurrentToken
		p.AdvanceToken()
		switch p.CurrentToken.Type {
		case token.MULTIPLY:
			node.Variant, err = parseWildcardImport(p)
		case token.LBRACE:
			node.Variant, err = parseNamedListImport(p)
		default:
			err = p.Error("syntax error")
		}
		if err != nil {
			return
		}
	} else {
		if node.Layout.From, err = p.ExpectLiteral("from"); err != nil {
			return
		}
		if node.From, err = p.Expect(token.STRING); err != nil {
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
	switch v := node.Variant.(type) {
	case *SideEffectImportStmt:
		printSideEffectImportStmt(pr, v)
	case *WildcardImportStmt:
		printWildcardImportStmt(pr, v)
	case *NamedListImportStmt:
		printNamedListImportStmt(pr, v)
	case *DefaultImportStmt:
		printDefaultImportStmt(pr, v)
	}
	return nil
}

func printSideEffectImportStmt(pr *printer.Printer, node *SideEffectImportStmt) {
	pr.Space().Print(node.From, node.Layout.Semi)
}

func printWildcardImportStmt(pr *printer.Printer, node *WildcardImportStmt) {
	pr.Space().Print(node.Layout.Multiply)
	pr.Space().Print(node.Layout.As)
	pr.Space().Print(node.Name)
	pr.Space().Print(node.Layout.From)
	pr.Space().Print(node.From, node.Layout.Semi)
}

func printNamedListImportStmt(pr *printer.Printer, node *NamedListImportStmt) {
	pr.Space().Print(node.Layout.Lbrace)
	if len(node.Items) > 0 {
		pr.IncreaseIndent()
		for i, item := range node.Items {
			if i > 0 {
				pr.Print(',')
				pr.Space()
			}
			pr.Space().Print(item.Pattern)
			if item.Alias != nil {
				pr.Space().Print(item.Layout.As)
				pr.Space().Print(item.Alias)
			}
		}
		pr.DecreaseIndent()
		pr.Space()
	}
	pr.Print(node.Layout.Rbrace)
	pr.Space().Print(node.Layout.From)
	pr.Space().Print(node.From, node.Layout.Semi)
}

func printDefaultImportStmt(pr *printer.Printer, node *DefaultImportStmt) {
	pr.Space().Print(node.Name)
	if node.Variant != nil {
		pr.Print(node.Layout.Comma)
		switch v := node.Variant.(type) {
		case *WildcardImportStmt:
			printWildcardImportStmt(pr, v)
		case *NamedListImportStmt:
			printNamedListImportStmt(pr, v)
		}
	} else {
		pr.Space().Print(node.Layout.From)
		pr.Space().Print(node.From, node.Layout.Semi)
	}
}
