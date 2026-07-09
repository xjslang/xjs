package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Builder struct {
	stmtParsers   []func(*Parser, func() (ast.Stmt, error)) (ast.Stmt, error)
	exprParsers   []func(*Parser, func() (ast.Expr, error)) (ast.Expr, error)
	unaryParsers  []func(*Parser, func() (ast.Expr, error)) (ast.Expr, error)
	binaryParsers []func(*Parser, ast.Expr, func(ast.Expr) (ast.Expr, error)) (ast.Expr, error)
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) UseStmtParser(parser func(p *Parser, next func() (ast.Stmt, error)) (ast.Stmt, error)) *Builder {
	b.stmtParsers = append(b.stmtParsers, parser)
	return b
}

func (b *Builder) UseExprParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) *Builder {
	b.exprParsers = append(b.exprParsers, parser)
	return b
}

func (b *Builder) UseUnaryParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) *Builder {
	b.unaryParsers = append(b.unaryParsers, parser)
	return b
}

func (b *Builder) UseBinaryParser(parser func(p *Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error)) *Builder {
	b.binaryParsers = append(b.binaryParsers, parser)
	return b
}

func (b *Builder) Build(sc token.Scanner) *Parser {
	p := &Parser{}
	for _, stmt := range b.stmtParsers {
		p.useStmtParser(stmt)
	}
	for _, expr := range b.exprParsers {
		p.useExprParser(expr)
	}
	for _, unaryExpr := range b.unaryParsers {
		p.useUnaryParser(unaryExpr)
	}
	for _, binaryExpr := range b.binaryParsers {
		p.useBinaryParser(binaryExpr)
	}
	p.init(sc)
	return p
}
