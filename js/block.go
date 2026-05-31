package js

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Block struct {
	LbraceToken token.Token
	RbraceToken token.Token

	Stmts []ast.Node
}

func (node *Block) Type() string {
	return "Block"
}

func ParseBlock(p *parser.Parser) (_ *Block, err error) {
	node := &Block{}
	if node.LbraceToken, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	var errs []error
	for p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE {
		prevToken := p.CurrentToken
		stmt, err := p.ParseStmt()
		if err != nil {
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
	if node.RbraceToken, err = p.Expect(token.RBRACE); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return node, errors.Join(errs...)
	}
	return node, nil
}

func PrintBlock(p *printer.Printer, node *Block) {
	p.Print(node.LbraceToken)
	p.IncreaseIndent()
	for _, stmt := range node.Stmts {
		p.Print(stmt)
	}
	// RBRACE is a special token, since the "leading trivia"
	// must be printed "before" indentation level decreases
	p.PrintTrivia(node.RbraceToken.LeadingTrivia)
	p.DecreaseIndent()
	p.LnPrint(node.RbraceToken.Literal)
}
