package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type ForinStmt struct {
	ast.BaseStmt
	Layout struct {
		For    token.Token
		Lparen token.Token
		Var    token.Token
		In     token.Token
		Rparen token.Token
	}
	VarDecl bool
	Pattern ast.Node
	Value   ast.Expr
	Then    ast.Stmt
}

func ParseForinStmt(p *parser.Parser) (node *ForinStmt, err error) {
	node = &ForinStmt{}
	if node.Layout.For, err = p.Expect(js.FOR); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if typ := p.CurrentToken.Type; typ == js.LET || typ == CONST || typ == VAR {
		node.VarDecl = true
		node.Layout.Var = p.CurrentToken
		p.AdvanceToken()
	}
	if node.Pattern, err = js.ParseRightExpr(p, IN.Precedence()); err != nil {
		return
	}
	if node.Layout.In, err = p.ExpectLiteral("in"); err != nil {
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

func PrintForinStmt(pr *printer.Printer, node *ForinStmt) error {
	pr.Line().Print(node.Layout.For)
	pr.Space().Print(node.Layout.Lparen)
	pr.IncreaseIndent()
	if node.VarDecl {
		pr.Print(node.Layout.Var)
		pr.Space()
	}
	pr.Print(node.Pattern)
	pr.Space().Print(node.Layout.In)
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
