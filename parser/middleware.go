package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) UsePrefixOpParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
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

func (p *Parser) UseInfixOpParser(parser func(p *Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (ast.Node, error)) {
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

func defaultPrefixExprParser(p *Parser) (node ast.Node, err error) {
	op := p.CurrentToken
	if op.Type == token.LPAREN {
		return ParseParenExpr(p)
	}
	nodeExpr := &ast.UnaryExpr{Operator: p.CurrentToken}
	p.AdvanceToken()
	if nodeExpr.Value, err = ParseValue(p); err != nil {
		return
	}
	node = nodeExpr
	return
}

func defaultInfixExprParser(p *Parser, leftVal ast.Node) (node ast.Node, err error) {
	op := p.CurrentToken
	nodeExpr := &ast.BinaryExpr{
		LeftValue: leftVal,
		Operator:  op,
	}
	p.AdvanceToken()
	if nodeExpr.RightValue, err = ParseRightValue(p, op.Type.Precedence()); err != nil {
		return
	}
	node = nodeExpr
	return
}

func defaultStmtParser(p *Parser) (ast.Node, error) {
	return ParseExprStmt(p)
}

func defaultExprParser(p *Parser) (val ast.Node, err error) {
	if val, err = ParseValue(p); err != nil {
		return
	}
	typ := p.CurrentToken.Type
	for typ.IsInfixOp() {
		if val, err = p.infixExprParser(p, val); err != nil {
			return
		}
		typ = p.CurrentToken.Type
	}
	return
}
