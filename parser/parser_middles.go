package parser

import (
	"fmt"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Handler func(p *Parser) ast.Expression

func (p *Parser) UsePrefixHandler(hook token.Type, middleware func(p *Parser, next Handler) ast.Expression) {
	next, ok := p.prefixParseFns[hook]
	if !ok {
		panic(fmt.Sprintf("no existing handler for hook %v", hook))
	}
	p.prefixParseFns[hook] = func(p *Parser) ast.Expression {
		return middleware(p, next)
	}
}

func (p *Parser) UseParseStatement(middleware func(p *Parser, next func(p *Parser) ast.Statement) ast.Statement) {
	next := p.parseStatement
	p.parseStatement = func(p *Parser) ast.Statement {
		return middleware(p, next)
	}
}

func (p *Parser) UseParseExpressionStatement(middleware func(p *Parser, next func(p *Parser) *ast.ExpressionStatement) *ast.ExpressionStatement) {
	next := p.parseExpressionStatement
	p.parseExpressionStatement = func(p *Parser) *ast.ExpressionStatement {
		return middleware(p, next)
	}
}
