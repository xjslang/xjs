package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

type Builder struct {
	LexerBuilder   *lexer.Builder
	stmMiddlewares []func(p *Parser, next func() ast.Statement) ast.Statement
	expMiddlewares []func(p *Parser, next func() ast.Expression) ast.Expression
}

func NewBuilder(lb *lexer.Builder) *Builder {
	return &Builder{
		LexerBuilder:   lb,
		stmMiddlewares: []func(p *Parser, next func() ast.Statement) ast.Statement{},
		expMiddlewares: []func(p *Parser, next func() ast.Expression) ast.Expression{},
	}
}

func (pb *Builder) Install(plugin func(*Builder)) *Builder {
	plugin(pb)
	return pb
}

func (pb *Builder) UseStatementParser(middleware func(p *Parser, next func() ast.Statement) ast.Statement) *Builder {
	pb.stmMiddlewares = append(pb.stmMiddlewares, middleware)
	return pb
}

func (pb *Builder) UseExpressionParser(middleware func(p *Parser, next func() ast.Expression) ast.Expression) *Builder {
	pb.expMiddlewares = append(pb.expMiddlewares, middleware)
	return pb
}

func (pb *Builder) RegisterPrefixOperator(tokenType token.Type, createExpr func(right func() ast.Expression) ast.Expression) {
	pb.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
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

func (pb *Builder) RegisterInfixOperator(tokenType token.Type, precedence int, createExpr func(token token.Token, left ast.Expression, right func() ast.Expression) ast.Expression) {
	pb.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.PeekToken.Type != tokenType {
			return next()
		}
		left := p.ParsePrefixExpression()
		p.NextToken() // consume operator
		right := func() ast.Expression {
			p.NextToken() // move to right expression
			return p.ParseExpressionWithPrecedence(precedence)
		}
		expr := createExpr(p.CurrentToken, left, right)
		return p.ParseRemainingExpression(expr)
	})
}

func (pb *Builder) RegisterOperand(tokenType token.Type, createExpr func() ast.Expression) {
	pb.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type != tokenType {
			return next()
		}
		expr := createExpr()
		return p.ParseRemainingExpression(expr)
	})
}

func (pb *Builder) Build(input string) *Parser {
	l := pb.LexerBuilder.Build(input)
	return NewWithOptions(l, ParserOptions{
		StatementMiddlewares:  pb.stmMiddlewares,
		ExpressionMiddlewares: pb.expMiddlewares,
	})
}
