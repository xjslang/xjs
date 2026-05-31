package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
)

func (p *Parser) UsePrefixParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.prefixExprParser
	if next == nil {
		next = defaultPrefixExprParser
	}
	p.prefixExprParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func (p *Parser) UseInfixParser(parser func(p *Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (ast.Node, error)) {
	next := p.infixExprParser
	if next == nil {
		next = defaultInfixExprParser
	}
	p.infixExprParser = func(p *Parser, leftVal ast.Node) (ast.Node, error) {
		return parser(p, leftVal, func(leftVal ast.Node) (ast.Node, error) {
			return next(p, leftVal)
		})
	}
}

func (p *Parser) UseStmtParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.stmtParser
	if next == nil {
		next = defaultStmtParser
	}
	p.stmtParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func (p *Parser) UseExprParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.exprParser
	if next == nil {
		next = defaultExprParser
	}
	p.exprParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func defaultPrefixExprParser(p *Parser) (ast.Node, error) {
	return nil, errors.New("unknown prefix expression")
}

func defaultInfixExprParser(p *Parser, leftVal ast.Node) (ast.Node, error) {
	return nil, errors.New("unknown infix expression")
}

func defaultStmtParser(p *Parser) (ast.Node, error) {
	return nil, errors.New("unknown statement")
}

func defaultExprParser(p *Parser) (val ast.Node, err error) {
	return nil, errors.New("unknown expression")
}
