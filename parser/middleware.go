package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
)

func (p *Parser) UseStatementParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.statementParser
	if next == nil {
		next = defaultStatementParser
	}
	p.statementParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func defaultStatementParser(p *Parser) (ast.Node, error) {
	switch p.CurrentToken.Type {
	case scanner.LET:
		return ParseLet(p)
	case scanner.FUNCTION:
		return ParseFunction(p)
	default:
		msg := "Unknown statement"
		p.AddError(msg)
		p.AdvanceToken() // consume unrecognizable token
		return nil, errors.New(msg)
	}
}
