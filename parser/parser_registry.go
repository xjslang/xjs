package parser

import "github.com/xjslang/xjs/ast"

func (p *Parser) RegisterPrefixOperator(literal string, createExpr func(right func() ast.Expression) ast.Expression) {
	p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Literal == literal {
			right := func() ast.Expression {
				p.NextToken()
				return p.ParseExpression()
			}
			return createExpr(right)
		}
		return next()
	})
}

func (p *Parser) RegisterInfixOperator(literal string, precedence int, createExpr func(left ast.Expression, right func() ast.Expression) ast.Expression) {
	p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.PeekToken.Literal != literal {
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

func (p *Parser) RegisterOperand(literal string, createExpr func() ast.Expression) {
	p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Literal != literal {
			return next()
		}
		expr := createExpr()
		return p.ParseRemainingExpression(expr)
	})
}
