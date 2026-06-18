package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type LabelStmt struct {
	ast.BaseStmt
	Layout struct {
		Colon token.Token
	}
	Name *Ident
	Stmt ast.Stmt
}

func (node *LabelStmt) SelfClosing() bool {
	if v, ok := node.Stmt.(SelfClosingStmt); ok {
		return v.SelfClosing()
	}
	return false
}

func ParseLabelStmt(p *parser.Parser) (node *LabelStmt, err error) {
	node = &LabelStmt{}
	if node.Name, err = ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Colon, err = p.Expect(token.COLON); err != nil {
		return
	}
	if node.Stmt, err = p.ParseStmt(); err != nil {
		return
	}
	return node, nil
}

func PrintLabelStmt(p *printer.Printer, node *LabelStmt) {
	p.Print(node.Name, node.Layout.Colon)
	p.SpPrint(node.Stmt)
}
