package parser

import (
	"testing"

	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func TestSemicolonInsertion(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		scoped bool
	}{
		{"explicit", "let x = 100;", false},
		{"implicit eof", "let x = 100", false},
		{"implicit newline", "let x = 100\nlet y = 200", false},
		{"implicit rbrace", "let x = 100}", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &scanner.Scanner{}
			l.Init([]byte(test.input))
			p := &Parser{}
			p.Init(l)
			if test.scoped {
				p.EnterScope(blockScope)
			}
			result, err := ParseLetStmt(p)
			if err != nil {
				t.Fatal(err)
			}
			tok := result.SemiToken
			expected := token.Token{Type: token.SEMICOLON, Literal: ";"}
			if tok.Type != expected.Type {
				t.Errorf("Expected %v, got %v", expected.Type, tok)
			} else if tok.Literal != expected.Literal {
				t.Errorf("Expected %q, got %q", expected.Literal, tok.Literal)
			}
		})
	}
}

func TestAdvanceToStatementEnd(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		scoped   bool
		expected token.Token
	}{
		{"empty", "", false, token.Token{Type: token.EOF}},
		{"semicolon starting", ";", false, token.Token{Type: token.EOF}},
		{"semicolon between", "let x = 100; let y = 200", false, token.Token{Type: token.LET}},
		{"newline between", "let x = 100\n let y = 200", false, token.Token{Type: token.LET}},
		{"end of line", "let x = 100", false, token.Token{Type: token.EOF}},
	}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &scanner.Scanner{}
			l.Init([]byte(test.input))
			p := &Parser{}
			p.Init(l)
			if test.scoped {
				p.EnterScope(blockScope)
			}
			p.AdvanceToStatementEnd()
			if got := p.CurrentToken.Type; got != test.expected.Type {
				t.Errorf("%d: Expected %v, got %v", i, test.expected.Type, got)
			}
		})
	}
}
