package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FOR = token.RegisterType("for")

type ForStmt struct {
	ast.BaseStmt
	Layout struct {
		For    token.Token
		Lparen token.Token
		Semi1  token.Token
		Semi2  token.Token
		Rparen token.Token
	}
	Init  ast.Stmt
	Cond  ast.Expr
	After ast.Expr
	Then  ast.Stmt
}

func ParseForStmt(p *parser.Parser) (node *ForStmt, err error) {
	node = &ForStmt{}
	if node.Layout.For, err = p.Expect(FOR); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if p.CurrentToken.Type != token.SEMICOLON {
		if node.Init, err = p.ParseStmt(); err != nil {
			return
		}
	} else {
		node.Layout.Semi1 = p.CurrentToken
		p.AdvanceToken()
	}
	if p.CurrentToken.Type != token.SEMICOLON {
		if node.Cond, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Semi2, err = p.Expect(token.SEMICOLON); err != nil {
		return
	}
	if p.CurrentToken.Type != token.RPAREN {
		if node.After, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	return node, nil
}

func PrintForStmt(p *printer.Printer, node *ForStmt) error {
	p.Line().Print(node.Layout.For)
	p.Space().Print(node.Layout.Lparen)
	p.IncreaseIndent()
	if node.Init != nil {
		p.Beside().Print(node.Init)
	}
	p.Print(node.Layout.Semi1)
	if node.Cond != nil {
		p.Space().Print(node.Cond)
	}
	p.Print(node.Layout.Semi2)
	if node.After != nil {
		p.Space().Print(node.After)
	}
	p.DecreaseIndent()
	p.Print(node.Layout.Rparen)
	switch v := node.Then.(type) {
	case *SemiStmt:
		p.Beside().Print(v)
	default:
		p.Space().Print(node.Then)
	}
	return nil
}
