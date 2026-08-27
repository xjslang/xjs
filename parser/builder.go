package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
)

type Builder struct {
	stmtParsers     []func(*Parser, func(*Parser) (ast.Stmt, error)) (ast.Stmt, error)
	exprParsers     []func(*Parser, func(*Parser) (ast.Expr, error)) (ast.Expr, error)
	prefixOpParsers []func(*Parser, func(*Parser) (ast.Expr, error)) (ast.Expr, error)
	infixOpParsers  []func(*Parser, ast.Expr, func(*Parser, ast.Expr) (ast.Expr, error)) (ast.Expr, error)
}

func (b *Builder) UseStmtParser(parser func(p *Parser, next func(*Parser) (ast.Stmt, error)) (ast.Stmt, error)) *Builder {
	b.stmtParsers = append(b.stmtParsers, parser)
	return b
}

func (b *Builder) UseExprParser(parser func(p *Parser, next func(*Parser) (ast.Expr, error)) (ast.Expr, error)) *Builder {
	b.exprParsers = append(b.exprParsers, parser)
	return b
}

func (b *Builder) UsePrefixOpParser(parser func(p *Parser, next func(*Parser) (ast.Expr, error)) (ast.Expr, error)) *Builder {
	b.prefixOpParsers = append(b.prefixOpParsers, parser)
	return b
}

func (b *Builder) UseInfixOpParser(parser func(p *Parser, left ast.Expr, next func(*Parser, ast.Expr) (ast.Expr, error)) (ast.Expr, error)) *Builder {
	b.infixOpParsers = append(b.infixOpParsers, parser)
	return b
}

func (b *Builder) Build(s *scanner.Scanner) *Parser {
	p := &Parser{}
	for _, stmt := range b.stmtParsers {
		p.useStmtParser(stmt)
	}
	for _, expr := range b.exprParsers {
		p.useExprParser(expr)
	}
	for _, prefixOpParser := range b.prefixOpParsers {
		p.usePrefixOpParser(prefixOpParser)
	}
	for _, infirOpParser := range b.infixOpParsers {
		p.useInfixOpParser(infirOpParser)
	}
	p.init(s)
	return p
}
