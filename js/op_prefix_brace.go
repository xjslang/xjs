package js

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type PrefixBraceOp struct {
	ast.BaseExpr
	Layout struct {
		Lbrace token.Token
		Rbrace token.Token
	}
	Entries []ast.Node
}

type ObjEntry struct {
	ast.BaseNode
	Key     ast.Node
	Value   ast.Expr
	Default ast.Expr
}

type ObjMethod struct {
	ast.BaseNode
	Layout struct {
		Flag     token.Token // get or set
		Multiply token.Token
	}
	isAccessor  bool
	IsGenerator bool
	Name        ast.Node
	Params      *FunctionParams
	Body        *BlockStmt
}

type ComputedExpr struct {
	ast.BaseExpr
	Layout struct {
		Lbracket token.Token
		Rbracket token.Token
	}
	Expr ast.Expr
}

func ParsePrefixBraceOp(p *parser.Parser) (node *PrefixBraceOp, err error) {
	node = &PrefixBraceOp{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RBRACE {
		var entry ast.Node
		if entry, err = parser.Switch(p,
			func(p *parser.Parser) (ast.Node, error) { return parseObjMethod(p) },
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

func parseObjMethod(p *parser.Parser) (node *ObjMethod, err error) {
	node = &ObjMethod{}
	if lit := p.CurrentToken.Literal; (lit == "get" || lit == "set") && p.PeekToken.Type != token.LPAREN {
		node.isAccessor = true
		node.Layout.Flag = p.CurrentToken
		p.AdvanceToken()
	}
	if node.IsGenerator = p.CurrentToken.Type == token.MULTIPLY; node.IsGenerator {
		node.Layout.Multiply = p.CurrentToken
		p.AdvanceToken()
	}
	switch p.CurrentToken.Type {
	case token.LBRACKET:
		if node.Name, err = ParseComputedExpr(p); err != nil {
			return
		}
	case token.STRING, token.NUMBER, SPREAD:
		if node.Name, err = ParseValue(p); err != nil {
			return
		}
	default:
		if node.Name, err = ParseObjKey(p); err != nil {
			return
		}
	}
	if node.Params, err = ParseFunctionParams(p); err != nil {
		return
	}
	if node.Body, err = ParseBlockStmt(p); err != nil {
		return
	}
	return
}

func parseObjEntry(p *parser.Parser) (node *ObjEntry, err error) {
	node = &ObjEntry{}
	switch p.CurrentToken.Type {
	case token.LBRACKET:
		if node.Key, err = ParseComputedExpr(p); err != nil {
			return
		}
	case token.STRING, token.NUMBER, SPREAD:
		if node.Key, err = ParseValue(p); err != nil {
			return
		}
	default:
		if node.Key, err = ParseObjKey(p); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.COLON {
		p.AdvanceToken()
		if node.Value, err = ParseRightExpr(p, token.ASSIGN.Precedence()); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.ASSIGN {
		p.AdvanceToken()
		if node.Default, err = ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
			return
		}
	}
	return
}

func PrintPrefixBraceOp(pr *printer.Printer, node *PrefixBraceOp) (err error) {
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
				case *ComputedExpr:
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
			case *ObjMethod:
				if v.isAccessor {
					pr.Print(v.Layout.Flag)
				}
				if v.IsGenerator {
					pr.Space().Print(v.Layout.Multiply)
					switch w := v.Name.(type) {
					case *ComputedExpr:
						pr.Print(w.Layout.Lbracket, w.Expr, w.Layout.Rbracket)
					default:
						pr.Print(w)
					}
				} else {
					pr.Space()
					switch w := v.Name.(type) {
					case *ComputedExpr:
						pr.Space().Print(w.Layout.Lbracket, w.Expr, w.Layout.Rbracket)
					default:
						pr.Space().Print(w)
					}
				}
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

func ParseObjKey(p *parser.Parser) (node *Ident, err error) {
	tok := p.CurrentToken
	if r, s := utf8.DecodeRuneInString(tok.Literal); s == 0 || !scanner.IsLetter(r) {
		err = p.Error("key expected")
		return
	}
	node = &Ident{Token: tok}
	p.AdvanceToken()
	return
}

func ParseComputedExpr(p *parser.Parser) (node *ComputedExpr, err error) {
	node = &ComputedExpr{}
	node.Layout.Lbracket = p.CurrentToken
	p.AdvanceToken()
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rbracket, err = p.Expect(token.RBRACKET); err != nil {
		return
	}
	return
}
