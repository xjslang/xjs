package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var DO = token.RegisterType("do")

type DoWhileStmt struct {
	ast.BaseStmt
	Layout struct {
		Do     token.Token
		While  token.Token
		Lparen token.Token
		Rparen token.Token
		Semi   token.Token
	}
	Cond ast.Expr
	Stmt ast.Stmt
}

func ParseDoWhileStmt(p *parser.Parser) (node *DoWhileStmt, err error) {
	node = &DoWhileStmt{}
	if node.Layout.Do, err = p.Expect(DO); err != nil {
		return
	}
	if node.Stmt, err = p.ParseStmt(); err != nil {
		return
	}
	if node.Layout.While, err = p.Expect(js.WHILE); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Cond, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Layout.Semi, err = js.ExpectSemi(p); err != nil {
		return
	}
	return
}

func PrintDoWhileStmt(pr *printer.Printer, node *DoWhileStmt) error {
	pr.Line().Print(node.Layout.Do)
	pr.Space().Print(node.Stmt)
	pr.Space().Print(node.Layout.While)
	pr.Space().Print(node.Layout.Lparen, node.Cond, node.Layout.Rparen)
	pr.Print(node.Layout.Semi)
	return nil
}
