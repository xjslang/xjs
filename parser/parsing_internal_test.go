package parser

import (
	"testing"

	"github.com/xjslang/xjs/scanner"
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
			result, err := ParseLet(p)
			if err != nil {
				t.Fatal(err)
			}
			tok := result.SemiToken
			expected := scanner.Token{Type: scanner.SEMICOLON, Literal: ";"}
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
		expected scanner.Token
	}{
		{"empty", "", false, scanner.Token{Type: scanner.EOF}},
		{"semicolon starting", ";", false, scanner.Token{Type: scanner.EOF}},
		{"semicolon between", "let x = 100; let y = 200", false, scanner.Token{Type: scanner.LET}},
		{"newline between", "let x = 100\n let y = 200", false, scanner.Token{Type: scanner.LET}},
		{"end of line", "let x = 100", false, scanner.Token{Type: scanner.EOF}},
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
			AdvanceToStatementEnd(p)
			if got := p.CurrentToken.Type; got != test.expected.Type {
				t.Errorf("%d: Expected %v, got %v", i, test.expected.Type, got)
			}
		})
	}
}
