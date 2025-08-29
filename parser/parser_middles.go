package parser

import (
	"fmt"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type PrefixParseFn func(p *Parser) ast.Expression
type InfixParseFn func(p *Parser, left ast.Expression) ast.Expression
type StatementParseFn func(p *Parser) ast.Statement

func (p *Parser) UsePrefixExpressionHandler(hook token.Type, middleware func(p *Parser, next PrefixParseFn) ast.Expression) {
	next, ok := p.prefixParseFns[hook]
	if !ok {
		// TODO: recover panic
		panic(fmt.Sprintf("no existing handler for hook %v", hook))
	}
	p.prefixParseFns[hook] = func(p *Parser) ast.Expression {
		return middleware(p, next)
	}
}

func (p *Parser) UseInfixExpressionHandler(hook token.Type, middleware func(p *Parser, left ast.Expression, next InfixParseFn) ast.Expression) {
	next, ok := p.infixParseFns[hook]
	if !ok {
		// TODO: recover panic
		panic(fmt.Sprintf("no existing handler for hook %v", hook))
	}
	p.infixParseFns[hook] = func(p *Parser, left ast.Expression) ast.Expression {
		return middleware(p, left, next)
	}
}

func (p *Parser) UseStatementHandler(middleware func(p *Parser, next StatementParseFn) ast.Statement) {
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

func (p *Parser) UseExpressionStatementHandler(middleware func(p *Parser, next func(p *Parser) *ast.ExpressionStatement) *ast.ExpressionStatement) {
	next := p.expressionStatementParseFn
	p.expressionStatementParseFn = func(p *Parser) *ast.ExpressionStatement {
		return middleware(p, next)
	}
}
