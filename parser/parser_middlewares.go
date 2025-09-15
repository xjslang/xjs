package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) useStatementInterceptor(interceptor func(p *Parser, next func() ast.Statement) ast.Statement) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return interceptor(p, func() ast.Statement {
			return next(p)
		})
	}
}

func (p *Parser) UseExpressionInterceptor(interceptor func(p *Parser, next func() ast.Expression) ast.Expression) {
	next := p.expressionParseFn
	p.expressionParseFn = func(p *Parser, precedence int) ast.Expression {
		oldPrecedence := p.currentExpressionPrecedence
		p.currentExpressionPrecedence = precedence
		defer func() {
			p.currentExpressionPrecedence = oldPrecedence
		}()
		return interceptor(p, func() ast.Expression {
			return next(p, precedence)
		})
	}
}
