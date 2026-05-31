package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func ParseExprStmt(p *Parser) (node *ast.ExprStmt, err error) {
	node = &ast.ExprStmt{}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return
}

func ParseRightValue(p *Parser, precedence int) (val ast.Node, err error) {
	if val, err = ParseValue(p); err != nil {
		return
	}
	for {
		typ := p.CurrentToken.Type
		if !typ.IsInfixOp() || precedence >= typ.Precedence() {
			break
		}
		if val, err = p.infixExprParser(p, val); err != nil {
			return
		}
	}
	return
}

func ParseValue(p *Parser) (ast.Node, error) {
	typ := p.CurrentToken.Type
	if typ.IsPrefixOp() {
		return p.prefixExprParser(p)
	}
	switch typ {
	case token.NUMBER, token.STRING, token.BOOLEAN:
		val := p.CurrentToken
		p.AdvanceToken()
		return &ast.BasicLit{Value: val}, nil
	case token.IDENT:
		val := p.CurrentToken
		p.AdvanceToken()
		return &ast.Ident{Value: val}, nil
	}
	msg := "Expected value"
	p.AddError(msg)
	return nil, errors.New(msg)
}
