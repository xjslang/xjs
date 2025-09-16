package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

type XJSBuilder struct {
	LexerBuilder     *lexer.Builder
	stmtInterceptors []func(p *XJSParser, next func() ast.Statement) ast.Statement
	expInterceptors  []func(p *XJSParser, next func() ast.Expression) ast.Expression
}

func NewBuilder(lb *lexer.Builder) *XJSBuilder {
	return &XJSBuilder{
		LexerBuilder:     lb,
		stmtInterceptors: []func(p *XJSParser, next func() ast.Statement) ast.Statement{},
		expInterceptors:  []func(p *XJSParser, next func() ast.Expression) ast.Expression{},
	}
}

func (pb *XJSBuilder) Install(plugin func(*XJSBuilder)) *XJSBuilder {
	plugin(pb)
	return pb
}

func (pb *XJSBuilder) UseStatementInterceptor(interceptor func(p *XJSParser, next func() ast.Statement) ast.Statement) *XJSBuilder {
	pb.stmtInterceptors = append(pb.stmtInterceptors, interceptor)
	return pb
}

func (pb *XJSBuilder) UseExpressionInterceptor(interceptor func(p *XJSParser, next func() ast.Expression) ast.Expression) *XJSBuilder {
	pb.expInterceptors = append(pb.expInterceptors, interceptor)
	return pb
}

func (pb *XJSBuilder) RegisterPrefixOperator(tokenType token.Type, createExpr func(token token.Token, right func() ast.Expression) ast.Expression) {
	pb.UseExpressionInterceptor(func(p *XJSParser, next func() ast.Expression) ast.Expression {
		if p.currentToken.Type != tokenType {
			return next()
		}
		right := func() ast.Expression {
			p.NextToken()
			return p.ParsePrefixExpression()
		}
		expr := createExpr(p.currentToken, right)
		return p.ParseRemainingExpression(expr)
	})
}

func (pb *XJSBuilder) RegisterInfixOperator(tokenType token.Type, precedence int, createExpr func(token token.Token, left ast.Expression, right func() ast.Expression) ast.Expression) {
	pb.UseExpressionInterceptor(func(p *XJSParser, next func() ast.Expression) ast.Expression {
		if p.peekToken.Type != tokenType {
			return next()
		}
		left := p.ParsePrefixExpression()
		p.NextToken() // consume operator
		right := func() ast.Expression {
			p.NextToken() // move to right expression
			return p.ParseExpressionWithPrecedence(precedence)
		}
		expr := createExpr(p.currentToken, left, right)
		return p.ParseRemainingExpression(expr)
	})
}

func (pb *XJSBuilder) RegisterOperand(tokenType token.Type, createExpr func(token token.Token) ast.Expression) {
	pb.RegisterPrefixOperator(tokenType, func(token token.Token, right func() ast.Expression) ast.Expression {
		return createExpr(token)
	})
}

func (pb *XJSBuilder) Build(input string) *XJSParser {
	l := pb.LexerBuilder.Build(input)
	return newWithOptions(l, parserOptions{
		stmtInterceptors: pb.stmtInterceptors,
		expInterceptors:  pb.expInterceptors,
	})
}
