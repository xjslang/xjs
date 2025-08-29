package parser

import (
	"fmt"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Handler func(p *Parser) ast.Expression
type InfixHandler func(p *Parser, left ast.Expression) ast.Expression
type StatementHandler func(p *Parser) ast.Statement

func (p *Parser) UsePrefixExpressionHandler(hook token.Type, middleware func(p *Parser, next Handler) ast.Expression) {
	next, ok := p.prefixParseFns[hook]
	if !ok {
		// TODO: recover panic
		panic(fmt.Sprintf("no existing handler for hook %v", hook))
	}
	p.prefixParseFns[hook] = func(p *Parser) ast.Expression {
		return middleware(p, next)
	}
}

func (p *Parser) UseInfixExpressionHandler(hook token.Type, middleware func(p *Parser, left ast.Expression, next InfixHandler) ast.Expression) {
	next, ok := p.infixParseFns[hook]
	if !ok {
		// TODO: recover panic
		panic(fmt.Sprintf("no existing handler for hook %v", hook))
	}
	p.infixParseFns[hook] = func(p *Parser, left ast.Expression) ast.Expression {
		return middleware(p, left, next)
	}
}

func (p *Parser) UseStatementHandler(middleware func(p *Parser, next StatementHandler) ast.Statement) {
	next := p.statementParseFn
	p.statementParseFn = func(p *Parser) ast.Statement {
		return middleware(p, next)
	}
}

func (p *Parser) UseParseExpressionStatement(middleware func(p *Parser, next func(p *Parser) *ast.ExpressionStatement) *ast.ExpressionStatement) {
	next := p.parseExpressionStatement
	p.parseExpressionStatement = func(p *Parser) *ast.ExpressionStatement {
		return middleware(p, next)
	}
}
