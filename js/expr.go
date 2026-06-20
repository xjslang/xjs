package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

type Variable struct {
	ast.BaseExpr
	Name token.Token
}

type Literal struct {
	ast.BaseExpr
	Value token.Token
}

func ParseExpr(p *parser.Parser) (val ast.Expr, err error) {
	if val, err = ParseValue(p); err != nil {
		return
	}
	typ := p.CurrentToken.Type
	for typ.IsBinaryOp() && !p.CurrentToken.AfterNewline {
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
		if !typ.IsBinaryOp() || p.CurrentToken.AfterNewline || precedence >= typ.Precedence() {
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
	case token.IDENT:
		val := p.CurrentToken
		p.AdvanceToken()
		return &Variable{Name: val}, nil
	case NUMBER, STRING:
		val := p.CurrentToken
		p.AdvanceToken()
		return &Literal{Value: val}, nil
	}
	return nil, parser.NewErrorAtToken(p.CurrentToken, "expression expected")
}
