package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type IndexExpr struct {
	ast.BaseExpr
	Layout struct {
		Lbracket token.Token
		Rbracket token.Token
	}
	Value ast.Expr
	Index ast.Expr
}

func ParseIndexExpr(p *parser.Parser, left ast.Expr) (node *IndexExpr, err error) {
	node = &IndexExpr{Value: left}
	if node.Layout.Lbracket, err = p.Expect(token.LBRACKET); err != nil {
		return
	}
	if node.Index, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rbracket, err = p.Expect(token.RBRACKET); err != nil {
		return nil, err
	}
	return node, nil
}

func PrintIndexExpr(p *printer.Printer, node *IndexExpr) error {
	p.Print(node.Value, node.Layout.Lbracket, node.Index, node.Layout.Rbracket)
	return nil
}
