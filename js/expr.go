package js

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func ParseExpr(p *parser.Parser) (val ast.Expr, err error) {
	if val, err = ParseValue(p); err != nil {
		return
	}
	typ := p.CurrentToken.Type
	for typ.IsBinaryOp() {
		if val, err = p.ParseBinaryExpr(val); err != nil {
			return
		}
		typ = p.CurrentToken.Type
	}
	return
}

func ParseRightExpr(p *parser.Parser, precedence int) (val ast.Expr, err error) {
	if val, err = ParseValue(p); err != nil {
		return
	}
	for {
		typ := p.CurrentToken.Type
		if !typ.IsBinaryOp() || precedence >= typ.Precedence() {
			break
		}
		if val, err = p.ParseBinaryExpr(val); err != nil {
			return
		}
	}
	return
}

func ParseValue(p *parser.Parser) (ast.Expr, error) {
	typ := p.CurrentToken.Type
	if typ.IsUnaryOp() {
		return p.ParseUnaryExpr()
	}
	switch typ {
	case token.NUMBER, token.STRING, token.BOOLEAN:
		val := p.CurrentToken
		p.AdvanceToken()
		return &Literal{Value: val}, nil
	case token.IDENT:
		val := p.CurrentToken
		p.AdvanceToken()
		return &Ident{Name: val}, nil
	}
	msg := "Expected value"
	p.AddError(msg)
	return nil, errors.New(msg)
}
