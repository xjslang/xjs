package js

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Ident struct {
	Name token.Token
}

func (node *Ident) Type() string {
	return "Ident"
}

type BasicLit struct {
	Value token.Token
}

func (node *BasicLit) Type() string {
	return "BasicLit"
}

func ParseValue(p *parser.Parser) (ast.Node, error) {
	typ := p.CurrentToken.Type
	if typ.IsPrefixOp() {
		return p.ParsePrefixExpr()
	}
	switch typ {
	case token.NUMBER, token.STRING, token.BOOLEAN:
		val := p.CurrentToken
		p.AdvanceToken()
		return &BasicLit{Value: val}, nil
	case token.IDENT:
		val := p.CurrentToken
		p.AdvanceToken()
		return &Ident{Name: val}, nil
	}
	msg := "Expected value"
	p.AddError(msg)
	return nil, errors.New(msg)
}

func ParseRightValue(p *parser.Parser, precedence int) (val ast.Node, err error) {
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

func PrintIdent(p *printer.Printer, node *Ident) {
	p.Print(node.Name)
}

func PrintBasicLit(p *printer.Printer, node *BasicLit) {
	p.Print(node.Value)
}
