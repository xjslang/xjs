package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type FunctionDecl struct {
	ast.BaseDecl
	Layout struct {
		Function token.Token
		Lparen   token.Token
		Rparen   token.Token
	}
	Name   *js.Ident
	Params []ast.Node
	Body   *js.BlockStmt
}

func ParseFunctionDecl(p *parser.Parser) (node *FunctionDecl, err error) {
	node = &FunctionDecl{}
	if node.Layout.Function, err = p.Expect(js.FUNCTION); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RPAREN {
		var param ast.Node
		switch p.CurrentToken.Type {
		case token.LBRACE:
			if param, err = ParseObjExpr(p); err != nil {
				return
			}
		case token.LBRACKET:
			if param, err = ParseArrayExpr(p); err != nil {
				return
			}
		default:
			if param, err = js.ParseIdent(p); err != nil {
				return
			}
		}
		node.Params = append(node.Params, param)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
		p.AdvanceToken()
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Body, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	return node, nil
}

func PrintFunctionDecl(pr *printer.Printer, node *FunctionDecl) error {
	pr.Line().Print(node.Layout.Function)
	pr.Space().Print(node.Name)
	pr.Print(node.Layout.Lparen)
	pr.IncreaseIndent()
	for i, param := range node.Params {
		if i > 0 {
			pr.Print(",")
			pr.Space()
		}
		pr.Print(param)
	}
	pr.DecreaseIndent()
	pr.Print(node.Layout.Rparen)
	pr.Space().Print(node.Body)
	return nil
}
