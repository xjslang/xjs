package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
)

func (p *Parser) UseUnaryParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.unaryExprParser
	if next == nil {
		next = defaultUnaryExprParser
	}
	p.unaryExprParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func (p *Parser) UseBinaryParser(parser func(p *Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (ast.Node, error)) {
	next := p.binaryExprParser
	if next == nil {
		next = defaultBinaryExprParser
	}
	p.binaryExprParser = func(p *Parser, leftVal ast.Node) (ast.Node, error) {
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

func defaultUnaryExprParser(p *Parser) (ast.Node, error) {
	return nil, errors.New("unknown unary expression")
}

func defaultBinaryExprParser(p *Parser, leftVal ast.Node) (ast.Node, error) {
	return nil, errors.New("unknown binary expression")
}

func defaultStmtParser(p *Parser) (ast.Node, error) {
	return nil, errors.New("unknown statement")
}

func defaultExprParser(p *Parser) (val ast.Node, err error) {
	return nil, errors.New("unknown expression")
}
