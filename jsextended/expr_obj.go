package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ObjEntry struct {
	ast.BaseNode
	Key     ast.Node
	Value   ast.Expr
	Default ast.Expr
}

type ObjAccessor struct {
	ast.BaseNode
	Layout struct {
		Flag token.Token // get or set
	}
	Name   *js.Ident
	Params *FunctionParams
	Body   *js.BlockStmt
}

type ObjExpr struct {
	ast.BaseExpr
	Layout struct {
		Lbrace token.Token
		Rbrace token.Token
	}
	Entries []ast.Node
}

func ParseObjExpr(p *parser.Parser) (node *ObjExpr, err error) {
	node = &ObjExpr{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RBRACE {
		var entry ast.Node
		if entry, err = parser.Switch(p,
			func(p *parser.Parser) (ast.Node, error) { return parseObjAccessor(p) },
			func(p *parser.Parser) (ast.Node, error) { return parseObjEntry(p) },
		); err != nil {
			return
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

func parseObjAccessor(p *parser.Parser) (node *ObjAccessor, err error) {
	node = &ObjAccessor{}
	if lit := p.CurrentToken.Literal; lit != "get" && lit != "set" {
		err = p.Error("get/set expected")
		return
	}
	node.Layout.Flag = p.CurrentToken
	p.AdvanceToken()
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

func parseObjEntry(p *parser.Parser) (node *ObjEntry, err error) {
	node = &ObjEntry{}
	switch p.CurrentToken.Type {
	case token.LBRACKET:
		if node.Key, err = js.ParseComputedExpr(p); err != nil {
			return
		}
	case token.STRING, token.NUMBER, SPREAD:
		if node.Key, err = js.ParseValue(p); err != nil {
			return
		}
	default:
		if node.Key, err = js.ParseObjKey(p); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.COLON {
		p.AdvanceToken()
		if node.Value, err = js.ParseRightExpr(p, token.ASSIGN.Precedence()); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.ASSIGN {
		p.AdvanceToken()
		if node.Default, err = p.ParseExpr(); err != nil {
			return
		}
	}
	return
}

func PrintObjExpr(pr *printer.Printer, node *ObjExpr) (err error) {
	pr.Print(node.Layout.Lbrace)
	if len(node.Entries) > 0 {
		pr.IncreaseIndent()
		for i, entry := range node.Entries {
			if i > 0 {
				pr.Print(",")
			}
			switch v := entry.(type) {
			case *ObjEntry:
				switch w := v.Key.(type) {
				case *js.ComputedExpr:
					pr.Space().Print(w.Layout.Lbracket, w.Expr, w.Layout.Rbracket)
				default:
					pr.Space().Print(w)
				}
				if v.Value != nil {
					pr.Print(":")
					pr.Space().Print(v.Value)
				}
				if v.Default != nil {
					pr.Space().Print("=")
					pr.Space().Print(v.Default)
				}
			case *ObjAccessor:
				pr.Print(v.Layout.Flag)
				pr.Space().Print(v.Name)
				if err = PrintFunctionParams(pr, v.Params); err != nil {
					return
				}
				pr.Space().Print(v.Body)
			default:
				err = pr.Error("object entry expected")
				return
			}
		}
		pr.DecreaseIndent()
		pr.Space()
	}
	pr.Print(node.Layout.Rbrace)
	return
}
