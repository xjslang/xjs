package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var LET = token.RegisterType("let")

type LetStmt struct {
	ast.StmtNode
	LetToken    token.Token
	AssignToken token.Token
	SemiToken   token.Token

	Name  token.Token
	Value ast.Expr
}

func ParseLetStmt(p *parser.Parser) (_ *LetStmt, err error) {
	node := &LetStmt{}
	if node.LetToken, err = p.Expect(LET); err != nil {
		return
	}
	if node.Name, err = p.Expect(token.IDENT); err != nil {
		return
	}
	if node.AssignToken, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return node, nil
}

func PrintLetStmt(p *printer.Printer, node *LetStmt) {
	p.LnPrint(node.LetToken)
	p.SpPrint(node.Name)
	p.SpPrint(node.AssignToken)
	p.SpPrint(node.Value)
	p.Print(node.SemiToken)
}
