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

func ParseBlockStmt(p *parser.Parser) (node *BlockStmt, err error) {
	node = &BlockStmt{}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	var errList parser.ErrorList
	for p.CurrentToken.Type != token.RBRACE && p.CurrentToken.Type != token.EOF {
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
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		errList = append(errList, err)
	}
	if errList != nil {
		return node, errList
	}
	return
}

func PrintBlockStmt(p *printer.Printer, node *BlockStmt) error {
	p.Print(node.Layout.Lbrace)
	if len(node.Stmts) > 0 {
		p.IncreaseIndent()
		var stmt ast.Stmt
		for _, stmt = range node.Stmts {
			p.Print(stmt)
		}
		p.DecreaseIndent()
		p.Line()
	}
	p.Print(node.Layout.Rbrace)
	return nil
}

func advanceToStmtEnd(p *parser.Parser) {
	for {
		typ := p.CurrentToken.Type
		if typ == token.SEMICOLON {
			p.AdvanceToken()
			break
		}
		if typ == token.EOF || typ == token.RBRACE || typ == token.LBRACE || p.CurrentToken.AfterNewline {
			break
		}
		p.AdvanceToken()
	}
}
