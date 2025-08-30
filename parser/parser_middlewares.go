package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) UseStatementHandler(handler func(p *Parser, next func(p *Parser) ast.Statement) ast.Statement) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return handler(p, next)
	}
}

func (p *Parser) UseExpressionHandler(handler func(p *Parser, precedence int, next func(*Parser, int) ast.Expression) ast.Expression) {
	next := p.expressionParseFn
	p.expressionParseFn = func(p *Parser, precedence int) ast.Expression {
		return handler(p, precedence, next)
	}
}
