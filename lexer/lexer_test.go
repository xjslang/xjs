package lexer

import (
	"strings"
	"testing"

	"github.com/xjslang/xjs/token"
)

func expectTokenSequence(t *testing.T, input string, expectedToks []token.Token) {
	l := New(strings.NewReader(input))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		if tok.Type != expectedTok.Type {
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		} else if tok.Literal != expectedTok.Literal {
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		}
	}
	tok := l.NextToken()
	if tok.Type != token.EOF || tok.Literal != "" {
		t.Errorf("Expected %v, got %q", token.EOF, tok.Literal)
	}
}

func TestScanContinuesAfterNullCharacter(t *testing.T) {
	expectTokenSequence(t, "Hello\x00World", []token.Token{
		{Type: token.IDENT, Literal: "Hello"},
		{Type: token.UNKNOWN, Literal: "\x00"},
		{Type: token.IDENT, Literal: "World"},
	})
}

func TestSkipWhitespaces(t *testing.T) {
	expectTokenSequence(t, "  one\ntwo\rthree\tfour ", []token.Token{
		{Type: token.IDENT, Literal: "one"},
		{Type: token.IDENT, Literal: "two"},
		{Type: token.IDENT, Literal: "three"},
		{Type: token.IDENT, Literal: "four"},
	})
}

func TestReadIden(t *testing.T) {
	expectTokenSequence(t, " hello  hello123   _hello123 ", []token.Token{
		{Type: token.IDENT, Literal: "hello"},
		{Type: token.IDENT, Literal: "hello123"},
		{Type: token.IDENT, Literal: "_hello123"},
	})
}

func TestReadNumber(t *testing.T) {
	expectTokenSequence(t, "123", []token.Token{{Type: token.NUMBER, Literal: "123"}})
}

func TestReadString(t *testing.T) {
	t.Run("legal string", func(t *testing.T) {
		expectTokenSequence(t, " 'Hello, World!' \"Hello, World!\"", []token.Token{
			{Type: token.STRING, Literal: "'Hello, World!'"},
			{Type: token.STRING, Literal: "\"Hello, World!\""},
		})
	})
	t.Run("illegal string", func(t *testing.T) {
		items := []string{
			"'Hello, World",  // missing '
			"'",              // missing '
			"\"Hello, World", // missing "
			"\"",             // missing "
		}
		expectTokenSequence(t, strings.Join(items, "\n"), []token.Token{
			{Type: token.ILLEGAL, Literal: "'Hello, World"},
			{Type: token.ILLEGAL, Literal: "'"},
			{Type: token.ILLEGAL, Literal: "\"Hello, World"},
			{Type: token.ILLEGAL, Literal: "\""},
		})
	})
}
