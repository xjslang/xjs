package js

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type BlockStmt struct {
	ast.BaseStmt
	Layout struct {
		Lbrace token.Token
		Rbrace token.Token
	}
	Stmts []ast.Stmt
}

func ParseBlockStmt(p *parser.Parser) (_ *BlockStmt, err error) {
	node := &BlockStmt{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	var errs []error
	var stmt ast.Stmt
	for i := 0; p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE; i++ {
		if i > 0 && !selfClosingStmt(stmt) {
			if stmt, err = ParseSemiStmt(p); err != nil {
				return
			}
			node.Stmts = append(node.Stmts, stmt)
		}
		if p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE {
			prevToken := p.CurrentToken
			if stmt, err = p.ParseStmt(); err != nil {
				if prevToken.Position == p.CurrentToken.Position {
					// advance position to avoid infinite loop
					p.AdvanceToken()
				}
				errs = append(errs, err)
				p.AdvanceToStmtEnd()
				continue
			}
			node.Stmts = append(node.Stmts, stmt)
		}
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return node, errors.Join(errs...)
	}
	return node, nil
}

func PrintBlockStmt(p *printer.Printer, node *BlockStmt) {
	p.Print(node.Layout.Lbrace)
	if len(node.Stmts) > 0 {
		p.IncreaseIndent()
		var stmt ast.Stmt
		for _, stmt = range node.Stmts {
			p.Print(stmt)
		}
		if p.LastChar() != '}' && p.LastChar() != ';' {
			p.Print(';')
		}
		p.DecreaseIndent()
		p.EnsureLine()
	}
	p.Print(node.Layout.Rbrace)
}
