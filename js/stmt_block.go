package js

import (
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

func (node *BlockStmt) SelfClosing() bool {
	return true
}

func ParseBlockStmt(p *parser.Parser) (node *BlockStmt, err error) {
	node = &BlockStmt{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	var errList parser.ErrorList
	var stmt ast.Stmt
	for i := 0; p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE; i++ {
		selfClosing := false
		if v, ok := stmt.(SelfClosingStmt); ok {
			selfClosing = v.SelfClosing()
		}
		if i > 0 && !selfClosing {
			if stmt, err = ParseSemiStmt(p); err != nil {
				errList = append(errList, err)
				continue
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
				if eList, ok := err.(parser.ErrorList); ok {
					errList = append(errList, eList...)
				} else {
					errList = append(errList, err)
				}
				AdvanceToStmtEnd(p)
				continue
			}
			node.Stmts = append(node.Stmts, stmt)
		}
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		errList = append(errList, err)
	}
	if errList != nil {
		return node, errList
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
		selfClosing := false
		if v, ok := stmt.(SelfClosingStmt); ok {
			selfClosing = v.SelfClosing()
		}
		if !selfClosing {
			p.Print(';')
		}
		p.DecreaseIndent()
		p.EnsureLine()
	}
	p.Print(node.Layout.Rbrace)
}
