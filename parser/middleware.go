package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) UseBinExprParser(parser func(p *Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (ast.Node, error)) {
	next := p.binExprParser
	if next == nil {
		next = defaultBinExprParser
	}
	p.binExprParser = func(p *Parser, leftVal ast.Node) (ast.Node, error) {
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

func defaultBinExprParser(p *Parser, leftVal ast.Node) (node ast.Node, err error) {
	op := p.CurrentToken
	if op.Type == token.LPAREN {
		return ParseCallExpr(p, leftVal)
	}
	nodeExpr := &ast.BinaryExpr{
		LeftValue: leftVal,
		Operator:  op,
	}
	if nodeExpr.RightValue, err = ParseRemainingExpr(p); err != nil {
		return
	}
	node = nodeExpr
	return
}

func defaultStmtParser(p *Parser) (ast.Node, error) {
	switch p.CurrentToken.Type {
	case token.LET:
		return ParseLetStmt(p)
	case token.FUNCTION:
		return ParseFuncDecl(p)
	default:
		return ParseExprStmt(p)
	}
}

func defaultExprParser(p *Parser) (val ast.Node, err error) {
	if val, err = p.parseValue(); err != nil {
		return
	}
	typ := p.CurrentToken.Type
	for typ.IsBinaryOperator() {
		if val, err = p.binExprParser(p, val); err != nil {
			return
		}
		typ = p.CurrentToken.Type
	}
	return
}
