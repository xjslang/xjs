package parser

import "github.com/xjslang/xjs/ast"

func (p *Parser) UseParseStatement(middleware func(p *Parser, next func(p *Parser) ast.Statement) ast.Statement) {
	next := p.parseStatement
	p.parseStatement = func(p *Parser) ast.Statement {
		return middleware(p, next)
	}
}

func (p *Parser) UseParseLetStatement(middleware func(p *Parser, next func(p *Parser) *ast.LetStatement) *ast.LetStatement) {
	next := p.parseLetStatement
	p.parseLetStatement = func(p *Parser) *ast.LetStatement {
		return middleware(p, next)
	}
}
