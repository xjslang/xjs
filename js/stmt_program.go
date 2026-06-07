package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Program struct {
	ast.BaseStmt
	Layout struct {
		EOF token.Token
	}
	Stmts []ast.Stmt
}

func ParseProgram(p *parser.Parser) (_ *Program, err error) {
	node := &Program{}
	var stmt ast.Stmt
	for i := 0; p.CurrentToken.Type != token.EOF; i++ {
		selfClosing := false
		if v, ok := stmt.(SelfClosingStmt); ok {
			selfClosing = v.SelfClosing()
		}
		if i > 0 && !selfClosing {
			if stmt, err = ParseSemiStmt(p); err != nil {
				return
			}
			node.Stmts = append(node.Stmts, stmt)
		}
		if p.CurrentToken.Type != token.EOF {
			prevToken := p.CurrentToken
			if stmt, err = p.ParseStmt(); err != nil {
				if prevToken.Position == p.CurrentToken.Position {
					// advance position to avoid infinite loop
					p.AdvanceToken()
				}
				p.AdvanceToStmtEnd()
				continue
			}
			node.Stmts = append(node.Stmts, stmt)
		}
	}
	node.Layout.EOF = p.CurrentToken
	if errors := p.Errors(); len(errors) > 0 {
		return node, errors
	}
	return node, nil
}

func PrintProgram(p *printer.Printer, node *Program) {
	var stmt ast.Stmt
	for _, stmt = range node.Stmts {
		p.Print(stmt)
	}
	selfClosing := false
	if v, ok := stmt.(SelfClosingStmt); ok {
		selfClosing = v.SelfClosing()
	}
	if len(node.Stmts) > 0 && !selfClosing {
		p.Print(';')
	}
	p.Print(node.Layout.EOF)
}
