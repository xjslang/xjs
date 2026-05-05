package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
)

func (p *Parser) UseStatementParser(parser func(p *Parser, next func() (ast.Statement, error)) (ast.Statement, error)) {
	next := p.statementParser
	if next == nil {
		next = defaultStatementParser
	}
	p.statementParser = func(p *Parser) (ast.Statement, error) {
		return parser(p, func() (ast.Statement, error) {
			return next(p)
		})
	}
}

func defaultStatementParser(p *Parser) (ast.Statement, error) {
	msg := "Unknown statement"
	p.AddError(msg)
	p.AdvanceToken() // consume unrecognizable token
	return nil, errors.New(msg)
}
