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
		} else if expectedTok.LeadingTrivia != nil && len(tok.LeadingTrivia) != len(expectedTok.LeadingTrivia) {
			t.Errorf("token %d: expected %d leading trivia lines, got %d", i, len(expectedTok.LeadingTrivia), len(tok.LeadingTrivia))
		} else {
			for j, line := range expectedTok.LeadingTrivia {
				if tok.LeadingTrivia[j] != line {
					t.Errorf("token %d: expected %q leading trivia line, got %q", i, line, tok.LeadingTrivia[j])
				}
			}
		}
	}
}

func TestEmptySinglelineComment(t *testing.T) {
	expectTokenSequence(t, "//\nhello//\r\nthere//\rObi-Wan Kenobi", []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{""}},
		{Type: token.IDENT, Literal: "there", LeadingTrivia: []string{""}},
		{Type: token.EOF, LeadingTrivia: []string{"\rObi-Wan Kenobi"}},
	})
}

func TestMultilineComments(t *testing.T) {
	expectTokenSequence(t, `/* lorem
ipsum dolor */

hello/* unfinished comment`, []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{" lorem\nipsum dolor ", "", ""}},
		{Type: token.ILLEGAL, Literal: " unfinished comment"},
		{Type: token.EOF},
	})
}

func TestSinglelineComments(t *testing.T) {
	expectTokenSequence(t, `
  // First Name
  John
  
  // Last Name
  Smith
	
	// Final comment`, []token.Token{
		{Type: token.IDENT, Literal: "John", LeadingTrivia: []string{"", " First Name"}},
		{Type: token.IDENT, Literal: "Smith", LeadingTrivia: []string{"", "", " Last Name"}},
		{Type: token.EOF, LeadingTrivia: []string{"", "", " Final comment"}},
	})
}

func TestScanContinuesAfterNullCharacter(t *testing.T) {
	expectTokenSequence(t, "Hello\x00World", []token.Token{
		{Type: token.IDENT, Literal: "Hello"},
		{Type: token.UNKNOWN, Literal: "\x00"},
		{Type: token.IDENT, Literal: "World"},
		{Type: token.EOF},
	})
}

func TestPunctuators(t *testing.T) {
	expectTokenSequence(t, "; = == ! != < <= > >= () {} /", []token.Token{
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.EQ, Literal: "=="},
		{Type: token.NOT, Literal: "!"},
		{Type: token.NOT_EQ, Literal: "!="},
		{Type: token.LT, Literal: "<"},
		{Type: token.LTE, Literal: "<="},
		{Type: token.GT, Literal: ">"},
		{Type: token.GTE, Literal: ">="},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.DIVIDE, Literal: "/"},
		{Type: token.EOF},
	})
}

func TestSkipWhitespaces(t *testing.T) {
	expectTokenSequence(t, "  one\ntwo\rthree\tfour \r\n five ", []token.Token{
		{Type: token.IDENT, Literal: "one"},
		{Type: token.IDENT, Literal: "two"},
		{Type: token.IDENT, Literal: "three"},
		{Type: token.IDENT, Literal: "four"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.EOF},
	})
}

func TestReadIden(t *testing.T) {
	expectTokenSequence(t, " hello  hello123   _hello123 ", []token.Token{
		{Type: token.IDENT, Literal: "hello"},
		{Type: token.IDENT, Literal: "hello123"},
		{Type: token.IDENT, Literal: "_hello123"},
		{Type: token.EOF},
	})
}

func TestReadNumber(t *testing.T) {
	expectTokenSequence(t, "123", []token.Token{
		{Type: token.NUMBER, Literal: "123"},
		{Type: token.EOF},
	})
}

func TestReadString(t *testing.T) {
	t.Run("legal string", func(t *testing.T) {
		expectTokenSequence(t, " 'Hello, World!' \"Hello, World!\"", []token.Token{
			{Type: token.STRING, Literal: "'Hello, World!'"},
			{Type: token.STRING, Literal: "\"Hello, World!\""},
			{Type: token.EOF},
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
			{Type: token.EOF},
		})
	})
}
