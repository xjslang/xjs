package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func NewBuilder(lb *lexer.Builder) *Builder {
	return &Builder{
		LexerBuilder:     lb,
		stmtInterceptors: []Interceptor[ast.Statement]{},
		expInterceptors:  []Interceptor[ast.Expression]{},
	}
}

func (pb *Builder) Install(plugin func(*Builder)) *Builder {
	plugin(pb)
	return pb
}

func (pb *Builder) UseStatementInterceptor(interceptor Interceptor[ast.Statement]) *Builder {
	pb.stmtInterceptors = append(pb.stmtInterceptors, interceptor)
	return pb
}

func (pb *Builder) UseExpressionInterceptor(interceptor Interceptor[ast.Expression]) *Builder {
	pb.expInterceptors = append(pb.expInterceptors, interceptor)
	return pb
}

func (pb *Builder) RegisterPrefixOperator(tokenType token.Type, createExpr func(token token.Token, right func() ast.Expression) ast.Expression) {
	pb.UseExpressionInterceptor(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type != tokenType {
			return next()
		}
		right := func() ast.Expression {
			p.NextToken()
			return p.ParsePrefixExpression()
		}
		expr := createExpr(p.CurrentToken, right)
		return p.ParseRemainingExpression(expr)
	})
}

func (pb *Builder) RegisterInfixOperator(tokenType token.Type, precedence int, createExpr func(token token.Token, left ast.Expression, right func() ast.Expression) ast.Expression) {
	pb.UseExpressionInterceptor(func(p *Parser, next func() ast.Expression) ast.Expression {
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

func (pb *Builder) RegisterOperand(tokenType token.Type, createExpr func(token token.Token) ast.Expression) {
	pb.RegisterPrefixOperator(tokenType, func(token token.Token, right func() ast.Expression) ast.Expression {
		return createExpr(token)
	})
}

func (pb *Builder) Build(input string) *Parser {
	l := pb.LexerBuilder.Build(input)
	return newWithOptions(l, parserOptions{
		stmtInterceptors: pb.stmtInterceptors,
		expInterceptors:  pb.expInterceptors,
	})
}
