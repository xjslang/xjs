package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) RegisterPrefixOperator(tokenType token.Type, createExpr func(right func() ast.Expression) ast.Expression) {
	p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type != tokenType {
			return next()
		}
		right := func() ast.Expression {
			p.NextToken()
			return p.ParsePrefixExpression()
		}
		expr := createExpr(right)
		return p.ParseRemainingExpression(expr)
	})
}

func (p *Parser) RegisterInfixOperator(tokenType token.Type, precedence int, createExpr func(left ast.Expression, right func() ast.Expression) ast.Expression) {
	p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.PeekToken.Type != tokenType {
			return next()
		}
		left := p.ParsePrefixExpression()
		p.NextToken() // consume operator
		right := func() ast.Expression {
			p.NextToken() // move to right expression
			return p.ParseExpressionWithPrecedence(precedence)
		}
		expr := createExpr(left, right)
		return p.ParseRemainingExpression(expr)
	})
}

func (p *Parser) RegisterOperand(tokenType token.Type, createExpr func() ast.Expression) {
	p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type != tokenType {
			return next()
		}
		expr := createExpr()
		return p.ParseRemainingExpression(expr)
	})
}
