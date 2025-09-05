// Package lexer provides lexical analysis functionality for the XJS language.
// It tokenizes source code into a sequence of tokens that can be consumed by the parser.
package lexer

import (
	"testing"

	"github.com/xjslang/xjs/token"
)

func TestRawStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "`hello there!`", "hello there!"},
		{"escaped_backtick", "`hello \\`there\\`!`", "hello `there`!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				if tok.Type != token.RAW_STRING || tok.Literal != tt.expected {
					t.Errorf("Lexer(%q) got %q, want %q", tt.input, tok.Literal, tt.expected)
					return
				}
			}
		})
	}
}
