package parser

import (
	"github.com/xjslang/xjs/ast"
)

// TODO: rename to UsseStatementParser(parser...) and update the README.md file
func (p *Parser) UseStatementHandler(handler func(p *Parser, next func() ast.Statement) ast.Statement) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return handler(p, func() ast.Statement {
			return next(p)
		})
	}
}

// TODO: rename to UseExpressionParser(parser...) and update the README.md file
func (p *Parser) UseExpressionHandler(handler func(p *Parser, next func() ast.Expression) ast.Expression) {
	next := p.expressionParseFn
	p.expressionParseFn = func(p *Parser, precedence int) ast.Expression {
		return handler(p, func() ast.Expression {
			return next(p, precedence)
		})
	}
}
