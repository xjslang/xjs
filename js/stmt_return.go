package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var RETURN = token.RegisterType("return")

type ReturnStmt struct {
	ast.BaseStmt
	Layout struct {
		Return token.Token
		Semi   token.Token
	}
	Value ast.Expr
}

func ParseReturnStmt(p *parser.Parser) (node *ReturnStmt, err error) {
	node = &ReturnStmt{}
	if node.Layout.Return, err = p.Expect(RETURN); err != nil {
		return
	}
	typ := p.CurrentToken.Type
	if !p.CurrentToken.AfterNewline && typ != token.EOF && typ != token.SEMICOLON && typ != token.RBRACE {
		if node.Value, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = ExpectSemi(p); err != nil {
		return
	}
	return node, nil
}

func PrintReturnStmt(p *printer.Printer, node *ReturnStmt) {
	p.Line().Print(node.Layout.Return)
	if node.Value != nil {
		p.Space().Print(node.Value)
	}
	p.Print(node.Layout.Semi)
}
