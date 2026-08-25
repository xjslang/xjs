package js

import (
	"github.com/xjslang/xjs/ast"
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
	VarDecl bool
	Pattern ast.Node
	Value   ast.Expr
	Then    ast.Stmt
}

func ParseForofStmt(p *parser.Parser) (node *ForofStmt, err error) {
	node = &ForofStmt{}
	if node.Layout.For, err = p.Expect(FOR); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if typ := p.CurrentToken.Type; typ == LET || typ == CONST || typ == VAR {
		node.VarDecl = true
		node.Layout.Var = p.CurrentToken
		p.AdvanceToken()
	}
	if node.Pattern, err = ParseRightExpr(p, IN.Precedence()); err != nil {
		return
	}
	if node.Layout.Of, err = p.ExpectLiteral("of"); err != nil {
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
	if node.VarDecl {
		pr.Print(node.Layout.Var)
		pr.Space()
	}
	pr.Print(node.Pattern)
	pr.Space().Print(node.Layout.Of)
	pr.Space().Print(node.Value)
	pr.DecreaseIndent()
	pr.Print(node.Layout.Rparen)
	switch v := node.Then.(type) {
	case *SemiStmt:
		pr.Beside().Print(v)
	default:
		pr.Space().Print(node.Then)
	}
	return nil
}
