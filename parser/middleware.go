package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) useUnaryParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	next := p.unaryExprParser
	p.unaryExprParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, func() (ast.Expr, error) {
			return next(p)
		})
	}
}

func (p *Parser) useBinaryParser(parser func(p *Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error)) {
	next := p.binaryExprParser
	p.binaryExprParser = func(p *Parser, left ast.Expr) (ast.Expr, error) {
		return parser(p, left, func(left ast.Expr) (ast.Expr, error) {
			return next(p, left)
		})
	}
}

func (p *Parser) useStmtParser(parser func(p *Parser, next func() (ast.Stmt, error)) (ast.Stmt, error)) {
	next := p.stmtParser
	p.stmtParser = func(p *Parser) (ast.Stmt, error) {
		return parser(p, func() (ast.Stmt, error) {
			return next(p)
		})
	}
}

func (p *Parser) useExprParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	next := p.exprParser
	p.exprParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, func() (ast.Expr, error) {
			return next(p)
		})
	}
}
