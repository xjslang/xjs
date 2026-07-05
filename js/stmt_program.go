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

func ParseProgram(p *parser.Parser) (node *Program, err error) {
	node = &Program{}
	var errList parser.ErrorList
	for p.CurrentToken.Type != token.EOF {
		prevToken := p.CurrentToken
		var stmt ast.Stmt
		if stmt, err = p.ParseStmt(); err != nil {
			if eList, ok := err.(parser.ErrorList); ok {
				errList = append(errList, eList...)
			} else {
				errList = append(errList, err)
			}
			if prevToken.Position == p.CurrentToken.Position {
				// advance position to avoid infinite loop
				p.AdvanceToken()
			}
			advanceToStmtEnd(p)
			continue
		}
		node.Stmts = append(node.Stmts, stmt)
	}
	node.Layout.EOF = p.CurrentToken
	if errList != nil {
		return node, errList
	}
	return
}

func PrintProgram(p *printer.Printer, node *Program) {
	var stmt ast.Stmt
	for _, stmt = range node.Stmts {
		p.Print(stmt)
	}
	p.Print(node.Layout.EOF)
}
