package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ForofStmt struct {
	ast.BaseStmt
	Layout struct {
		For    token.Token
		Lparen token.Token
		Var    token.Token
		Of     token.Token
		Rparen token.Token
	}
	Pattern ast.Node
	Value   ast.Expr
	Then    ast.Stmt
}

func ParseForofStmt(p *parser.Parser) (node *ForofStmt, err error) {
	node = &ForofStmt{}
	if node.Layout.For, err = p.Expect(js.FOR); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if typ := p.CurrentToken.Type; typ != js.LET && typ != CONST && typ != VAR {
		err = p.Error("syntax error")
		return
	}
	node.Layout.Var = p.CurrentToken
	p.AdvanceToken()
	switch p.CurrentToken.Type {
	case token.LBRACE:
		node.Pattern, err = ParseObjExpr(p)
	case token.LBRACKET:
		node.Pattern, err = ParseArrayExpr(p)
	default:
		node.Pattern, err = js.ParseIdent(p)
	}
	if err != nil {
		return
	}
	if node.Layout.Of, err = p.ExpectString("of"); err != nil {
		return
	}
	if node.Value, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	return
}

func PrintForofStmt(pr *printer.Printer, node *ForofStmt) error {
	pr.Line().Print(node.Layout.For)
	pr.Space().Print(node.Layout.Lparen)
	pr.IncreaseIndent()
	pr.Print(node.Layout.Var)
	pr.Space().Print(node.Pattern)
	pr.Space().Print(node.Layout.Of)
	pr.Space().Print(node.Value)
	pr.DecreaseIndent()
	pr.Print(node.Layout.Rparen)
	switch v := node.Then.(type) {
	case *js.SemiStmt:
		pr.Beside().Print(v)
	default:
		pr.Space().Print(node.Then)
	}
	return nil
}
