package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) UseStatementHandler(middleware func(p *Parser, next func(p *Parser) ast.Statement) ast.Statement) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return middleware(p, next)
	}
}

func (p *Parser) UseExpressionHandler(middleware func(p *Parser, precedence int, next func(*Parser, int) ast.Expression) ast.Expression) {
	next := p.expressionParseFn
	p.expressionParseFn = func(p *Parser, precedence int) ast.Expression {
		return middleware(p, precedence, next)
	}
}
