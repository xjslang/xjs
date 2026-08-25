package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Ident struct {
	ast.BaseNode
	token.Token
}

type Variable struct {
	ast.BaseExpr
	token.Token
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
	for typ.IsInfixOp() {
		if val, err = p.ParseInfixOp(val); err != nil {
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
		if !typ.IsInfixOp() || precedence >= typ.Precedence() {
			break
		}
		if val, err = p.ParseInfixOp(val); err != nil {
			return
		}
	}
	return
}

func ParseValue(p *parser.Parser) (ast.Expr, error) {
	typ := p.CurrentToken.Type
	if typ.IsPrefixOp() {
		return p.ParsePrefixOp()
	}
	switch typ {
	case token.IDENT:
		val := p.CurrentToken
		p.AdvanceToken()
		return &Variable{Token: val}, nil
	case token.NUMBER, token.STRING:
		val := p.CurrentToken
		p.AdvanceToken()
		return &Literal{Value: val}, nil
	case token.UNCLOSED_STRING:
		return nil, p.Error("unclosed string")
	case token.UNCLOSED_BLOCK_COMMENT:
		return nil, p.Error("unclosed block comment")
	case token.INVALID_NUMBER:
		return nil, p.Error("invalid number")
	}
	return nil, p.Error("unknown expression")
}

func ParseIdent(p *parser.Parser) (node *Ident, err error) {
	node = &Ident{}
	if node.Token, err = p.Expect(token.IDENT); err != nil {
		return
	}
	return node, nil
}

func PrintIdent(pr *printer.Printer, node *Ident) error {
	pr.Print(node.Token)
	return nil
}

func IsSemi(tok token.Token) bool {
	if typ := tok.Type; typ == token.SEMICOLON || typ == token.RBRACE || typ == token.RPAREN || typ == token.EOF {
		return true
	}
	return tok.IsAfterNewline()
}
