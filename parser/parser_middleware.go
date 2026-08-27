package parser

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Parser) usePrefixOpParser(parser func(p *Parser, next func(*Parser) (ast.Expr, error)) (ast.Expr, error)) {
	next := p.prefixOpParser
	if next == nil {
		next = defaultPrefixOpParser
	}
	p.prefixOpParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, next)
	}
}

func (p *Parser) useInfixOpParser(parser func(p *Parser, left ast.Expr, next func(*Parser, ast.Expr) (ast.Expr, error)) (ast.Expr, error)) {
	next := p.infixOpParser
	if next == nil {
		next = infixOpParser
	}
	p.infixOpParser = func(p *Parser, left ast.Expr) (ast.Expr, error) {
		return parser(p, left, next)
	}
}

func (p *Parser) useStmtParser(parser func(p *Parser, next func(*Parser) (ast.Stmt, error)) (ast.Stmt, error)) {
	next := p.stmtParser
	if next == nil {
		next = defaultStmtParser
	}
	p.stmtParser = func(p *Parser) (ast.Stmt, error) {
		return parser(p, next)
	}
}

func (p *Parser) useExprParser(parser func(p *Parser, next func(*Parser) (ast.Expr, error)) (ast.Expr, error)) {
	next := p.exprParser
	if next == nil {
		next = defaultExprParser
	}
	p.exprParser = func(p *Parser) (ast.Expr, error) {
		return parser(p, next)
	}
}

func defaultPrefixOpParser(p *Parser) (ast.Expr, error) {
	return nil, p.Error("unknown prefix operator")
}

func infixOpParser(p *Parser, _ ast.Expr) (ast.Expr, error) {
	return nil, p.Error("unknown infix operator")
}

func defaultStmtParser(p *Parser) (ast.Stmt, error) {
	return nil, p.Error("unknown statement")
}

func defaultExprParser(p *Parser) (ast.Expr, error) {
	return nil, p.Error("unknown expression")
}
