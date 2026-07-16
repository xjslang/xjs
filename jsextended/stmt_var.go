package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	CONST = token.RegisterType("const")
	VAR   = token.RegisterType("var")
)

type VarStmt struct {
	ast.BaseDecl
	Layout struct {
		Var    token.Token
		Assign token.Token
		Semi   token.Token
	}
	Pattern ast.Node
	Value   ast.Expr
}

func ParseVarStmt(p *parser.Parser) (node *VarStmt, err error) {
	node = &VarStmt{}
	if typ := p.CurrentToken.Type; typ != js.LET && typ != CONST && typ != VAR {
		err = p.Error("syntax error")
		return
	}
	node.Layout.Var = p.CurrentToken
	p.AdvanceToken()
	switch p.CurrentToken.Type {
	case token.LBRACE:
		if node.Pattern, err = ParseObjExpr(p); err != nil {
			return
		}
	case token.LBRACKET:
		if node.Pattern, err = ParseArrayExpr(p); err != nil {
			return
		}
	default:
		if node.Pattern, err = js.ParseIdent(p); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.ASSIGN {
		node.Layout.Assign = p.CurrentToken
		p.AdvanceToken()
		if node.Value, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = js.ExpectSemi(p); err != nil {
		return
	}
	return
}

func PrintVarStmt(pr *printer.Printer, node *VarStmt) error {
	pr.Line().Print(node.Layout.Var)
	pr.Space().Print(node.Pattern)
	if node.Value != nil {
		pr.Space().Print(node.Layout.Assign)
		pr.Space().Print(node.Value)
	}
	pr.Print(node.Layout.Semi)
	return nil
}
