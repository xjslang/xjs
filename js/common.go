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
	for typ.IsBinaryOp() {
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
		if !typ.IsBinaryOp() || precedence >= typ.Precedence() {
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
	}
	return nil, p.Error("expression expected")
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

// ExpectSemi expects a semicolon or any other symbol that acts as a
// "statement terminator", such as '}' or ')'. If the statement terminator is a
// semicolon, then it consumes it and advances to the next token.
func ExpectSemi(p *parser.Parser) (tok token.Token, err error) {
	tok = p.CurrentToken
	switch tok.Type {
	case token.SEMICOLON:
		p.AdvanceToken()
		return
	case token.RBRACE, token.RPAREN, token.EOF:
		tok = token.Token{
			Type:     token.SEMICOLON,
			Literal:  token.SEMICOLON.Literal(),
			Position: tok.Position,
		}
		return
	default:
		if tok.AfterNewline {
			tok = token.Token{
				Type:     token.SEMICOLON,
				Literal:  token.SEMICOLON.Literal(),
				Position: tok.Position,
			}
			return
		}
	}
	err = p.Error(token.SEMICOLON.Literal() + " expected")
	return
}

func IsSemi(tok token.Token) bool {
	if typ := tok.Type; typ == token.SEMICOLON || typ == token.RBRACE || typ == token.RPAREN || typ == token.EOF {
		return true
	}
	return tok.AfterNewline
}
