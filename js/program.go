package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Program struct {
	EOFToken token.Token

	Stmts []ast.Node
}

func (node *Program) Type() string {
	return "Program"
}

func ParseProgram(p *parser.Parser) (_ *Program, err error) {
	node := &Program{}
	for p.CurrentToken.Type != token.EOF {
		prevToken := p.CurrentToken
		stmt, err := p.ParseExprStmt()
		if err != nil {
			if prevToken.Position == p.CurrentToken.Position {
				// advance position to avoid infinite loop
				p.AdvanceToken()
			}
			p.AdvanceToStmtEnd()
			continue
		}
		node.Stmts = append(node.Stmts, stmt)
	}
	node.EOFToken = p.CurrentToken
	if errors := p.Errors(); len(errors) > 0 {
		return node, errors
	}
	return node, nil
}

func PrintProgram(p *printer.Printer, node *Program) {
	for _, stmt := range node.Stmts {
		p.Print(stmt)
	}
	p.Print(node.EOFToken)
}
