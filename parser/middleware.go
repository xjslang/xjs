package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) UseUnaryParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	next := p.unaryExprParser
	p.unaryExprParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, func() (ast.Expr, error) {
			return next(p)
		})
	}
}

func (p *Parser) UseBinaryParser(parser func(p *Parser, leftVal ast.Expr, next func(leftVal ast.Expr) (ast.Expr, error)) (ast.Expr, error)) {
	next := p.binaryExprParser
	p.binaryExprParser = func(p *Parser, leftVal ast.Expr) (ast.Expr, error) {
		return parser(p, leftVal, func(leftVal ast.Expr) (ast.Expr, error) {
			return next(p, leftVal)
		})
	}
}

func (p *Parser) UseStmtParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.stmtParser
	p.stmtParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func (p *Parser) UseExprParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	next := p.exprParser
	p.exprParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, func() (ast.Expr, error) {
			return next(p)
		})
	}
}
