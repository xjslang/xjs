package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) useUnaryParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	next := p.unaryExprParser
	if next == nil {
		next = defaultUnaryParser
	}
	p.unaryExprParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, func() (ast.Expr, error) {
			return next(p)
		})
	}
}

func (p *Parser) useBinaryParser(parser func(p *Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error)) {
	next := p.binaryExprParser
	if next == nil {
		next = defaultBinaryParser
	}
	p.binaryExprParser = func(p *Parser, left ast.Expr) (ast.Expr, error) {
		return parser(p, left, func(left ast.Expr) (ast.Expr, error) {
			return next(p, left)
		})
	}
}

func (p *Parser) useStmtParser(parser func(p *Parser, next func() (ast.Stmt, error)) (ast.Stmt, error)) {
	next := p.stmtParser
	if next == nil {
		next = defaultStmtParser
	}
	p.stmtParser = func(p *Parser) (ast.Stmt, error) {
		return parser(p, func() (ast.Stmt, error) {
			return next(p)
		})
	}
}

func (p *Parser) useExprParser(parser func(p *Parser, next func() (ast.Expr, error)) (ast.Expr, error)) {
	next := p.exprParser
	if next == nil {
		next = defaultExprParser
	}
	p.exprParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, func() (ast.Expr, error) {
			return next(p)
		})
	}
}

func defaultUnaryParser(p *Parser) (ast.Expr, error) {
	return nil, p.Error("unknown unary operator")
}

func defaultBinaryParser(p *Parser, _ ast.Expr) (ast.Expr, error) {
	return nil, p.Error("unknown binary operator")
}

func defaultStmtParser(p *Parser) (ast.Stmt, error) {
	return nil, p.Error("unknown statement")
}

func defaultExprParser(p *Parser) (ast.Expr, error) {
	return nil, p.Error("unknown expression")
}
