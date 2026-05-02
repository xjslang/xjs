package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) UseStatementParser(parser func(p *Parser, next func() ast.Statement) ast.Statement) {
	next := p.statementParser
	if next == nil {
		next = defaultStatementParser
	}
	p.statementParser = func(p *Parser) ast.Statement {
		return parser(p, func() ast.Statement {
			return next(p)
		})
	}
}

func defaultStatementParser(p *Parser) ast.Statement {
	switch p.CurrentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.FUNCTION:
		return p.parseFunction()
	}
	return nil
}
