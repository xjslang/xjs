package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) UseStatementParser(parseFn func(p *Parser, next func() ast.Statement) ast.Statement) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return parseFn(p, func() ast.Statement {
			return next(p)
		})
	}
}

func (p *Parser) UseExpressionParser(parseFn func(p *Parser, next func(left ast.Expression) ast.Expression) ast.Expression) {
	originalExpressionParseFn := p.expressionParseFn

	p.expressionParseFn = func(p *Parser, precedence int) ast.Expression {
		return parseFn(p, func(left ast.Expression) ast.Expression {
			if left != nil {
				return p.parseRemainingExpression(left, precedence)
			}
			return originalExpressionParseFn(p, precedence)
		})
	}
}
