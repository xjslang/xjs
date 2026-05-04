package parser_test

import (
	"errors"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type DeferStatement struct {
	//
}

func (node *DeferStatement) PrintTo(p *printer.Printer) {
	//
}

func DeferPlugin(b *parser.Builder) {
	b.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "defer" {
			p.AdvanceToken() // consume defer
			if p.CurrentToken.Type != token.IDENT {
				msg := "Expected " + token.IDENT.String()
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			p.AdvanceToken()
			if p.CurrentToken.Type != token.LPAREN {
				msg := "Expected " + token.LPAREN.String()
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			p.AdvanceToken()
			if p.CurrentToken.Type != token.RPAREN {
				msg := "Expected " + token.RPAREN.String()
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			p.AdvanceToken()
			return &DeferStatement{}, nil
		}
		return next()
	})
}

func TestBuilder(t *testing.T) {
	l := &lexer.Lexer{}
	l.Init([]byte(`defer dbClose()`))

	b := parser.Builder{}
	b.Install(DeferPlugin)
	p := b.Build(l)
	pr, err := p.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}
	if n := len(pr.Statements); n != 1 {
		t.Fatalf("Expected 1 statements, got %d", n)
	}
	stmt := pr.Statements[0]
	if _, ok := stmt.(*DeferStatement); !ok {
		t.Fatalf("Expected *DeferStatement, got %T", stmt)
	}
}
