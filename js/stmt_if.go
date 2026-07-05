package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	IF   = token.RegisterType("if")
	ELSE = token.RegisterType("else")
)

type IfStmt struct {
	ast.BaseStmt
	Layout struct {
		If     token.Token
		Lparen token.Token
		Rparen token.Token
		Else   token.Token
	}
	Cond ast.Expr
	Then ast.Stmt
	Else ast.Stmt
}

func (node *IfStmt) SelfClosing() bool {
	if node.Else != nil {
		if v, ok := node.Else.(SelfClosingStmt); ok {
			return v.SelfClosing()
		}
		return false
	}
	if v, ok := node.Then.(SelfClosingStmt); ok {
		return v.SelfClosing()
	}
	return false
}

func ParseIfStmt(p *parser.Parser) (node *IfStmt, err error) {
	node = &IfStmt{}
	// if
	if node.Layout.If, err = p.Expect(IF); err != nil {
		return
	}
	// (condition)
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Cond, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	// then
	if p.CurrentToken.Type == token.LBRACE {
		if node.Then, err = p.ParseStmt(); err != nil {
			return
		}
		// else
		if p.CurrentToken.Type == ELSE {
			node.Layout.Else = p.CurrentToken
			p.AdvanceToken()
			if node.Else, err = p.ParseStmt(); err != nil {
				return
			}
		}
	} else if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	return
}

func PrintIfStmt(p *printer.Printer, node *IfStmt) {
	// if (condition) stmt
	p.Line().Print(node.Layout.If)
	p.Space().Print(node.Layout.Lparen)
	p.Print(node.Cond, node.Layout.Rparen)
	p.Space().Print(node.Then)
	// else
	if node.Else != nil {
		p.Space().Print(node.Layout.Else)
		p.Space().Print(node.Else)
	}
}
