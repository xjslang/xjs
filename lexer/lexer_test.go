package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/token"
)

type tokenCompareConfig struct {
	afterNewline  bool
	leadingTrivia bool
	tokenPosition bool
}

type tokenCompareOption func(cfg *tokenCompareConfig)

func compareAfterNewline() tokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.afterNewline = true
	}
}

func compareLeadingTrivia() tokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.leadingTrivia = true
	}
}

func compareTokenPosition() tokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.tokenPosition = true
	}
}

func assertTokens(t *testing.T, input string, expectedToks []token.Token, opts ...tokenCompareOption) {
	cfg := &tokenCompareConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	l := New([]byte(input))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		case cfg.afterNewline && tok.AfterNewline != expectedTok.AfterNewline:
			t.Errorf("token %d: expected AfterNewline to be %t, got %t", i, expectedTok.AfterNewline, tok.AfterNewline)
		case cfg.leadingTrivia:
			if len(tok.LeadingTrivia) != len(expectedTok.LeadingTrivia) {
				t.Errorf("token %d: expected %d leading trivia lines, got %d", i, len(expectedTok.LeadingTrivia), len(tok.LeadingTrivia))
			} else {
				for j, line := range expectedTok.LeadingTrivia {
					if tok.LeadingTrivia[j] != line {
						t.Errorf("token %d: expected %q leading trivia line, got %q", i, line, tok.LeadingTrivia[j])
					}
				}
			}
		case cfg.tokenPosition && (tok.Line != expectedTok.Line || tok.Column != expectedTok.Column):
			t.Errorf("token %d: expected position to be (%d, %d), got (%d, %d)", i, expectedTok.Line, expectedTok.Column, tok.Line, tok.Column)
		}
	}
}

func BenchmarkLexer(b *testing.B) {
	l := New([]byte("lorem ipsum dolor"))
	var tok token.Token // prevent dead code elimination
	for b.Loop() {
		for tok = l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		}
		l.Reset()
	}
	_ = tok
}

func TestTokenPosition(t *testing.T) {
	input := " aaa   bbb /* block comment*/ ccc\n// comment\rddd\r\ne!\n"
	assertTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "aaa", Line: 0, Column: 1},
		{Type: token.IDENT, Literal: "bbb", Line: 0, Column: 7},
		{Type: token.IDENT, Literal: "ccc", Line: 0, Column: 30},
		{Type: token.IDENT, Literal: "ddd", Line: 2, Column: 0},
		{Type: token.IDENT, Literal: "e", Line: 3, Column: 0},
		{Type: token.NOT, Literal: "!", Line: 3, Column: 1},
		{Type: token.EOF, Line: 4, Column: 0},
	}, compareTokenPosition())
}

func TestReset(t *testing.T) {
	items := []string{"lorem", "ipsum", "dolor"}
	l := New([]byte(strings.Join(items, " ")))
	for range 2 {
		var toks []token.Token
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			toks = append(toks, tok)
		}
		if l := len(toks); l != len(items) {
			t.Fatalf("Expected len(toks)=%d, got %d", len(items), l)
		}
		for i, tok := range toks {
			if tok.Type != token.IDENT {
				t.Fatalf("token %d: expected type %v, got %v", i, token.IDENT, tok.Type)
			}
			if tok.Literal != items[i] {
				t.Fatalf("token %d: expected %q, got %q", i, items[i], tok.Literal)
			}
		}
		l.Reset()
	}
}

func TestUnicodeChars(t *testing.T) {
	tests := []struct {
		name  string
		items []string
	}{
		{"general", []string{"España", "Türkiye", "São Tomé e Príncipe", "Curaçao", "Réunion"}},
		{"various diacritics", []string{"á", "é", "í", "ó", "ú", "ü", "ñ", "ç", "ø", "å", "ä", "ö"}},
		{"emojis", []string{"🇪🇸", "👍", "🎉"}},
		{"non-Latin alphabets", []string{"Россия", "مصر", "中国", "日本", "한국"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var expectedToks []token.Token
			for _, item := range test.items {
				expectedToks = append(expectedToks, token.Token{Type: token.STRING, Literal: fmt.Sprintf("'%s'", item)})
			}
			expectedToks = append(expectedToks, token.Token{Type: token.EOF})
			item := "'" + strings.Join(test.items, "' '") + "'"
			assertTokens(t, item, expectedToks)
		})
	}
}

func TestAfterNewline(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"newline before block comment", "hello\n/* block comment */world"},
		{"block comment with \n in the middle", "hello/* block\ncomment */world"},
		{"block comment with \r in the middle", "hello/* block\rcomment */world"},
		{"block comment with \r\n in the middle", "hello/* block\r\ncomment */world"},
		{"single-line comment", "hello// comment\nworld"},
		{"newline", "hello\nworld"},
		{"newline", "hello\rworld"},
		{"newline", "hello\r\nworld"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertTokens(t, test.input, []token.Token{
				{Type: token.IDENT, Literal: "hello"},
				{Type: token.IDENT, Literal: "world", AfterNewline: true},
			}, compareAfterNewline())
		})
	}
}

func TestBlockComments(t *testing.T) {
	input := `/* lorem
ipsum dolor */

hello/* unfinished comment`
	assertTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{" lorem\nipsum dolor ", "", ""}},
		{Type: token.ILLEGAL, Literal: " unfinished comment"},
		{Type: token.EOF},
	}, compareLeadingTrivia())
}

func TestLineComments(t *testing.T) {
	input := `
  // First Name
  John
  
  // Last Name
  Smith
	
	// Final comment`
	assertTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "John", LeadingTrivia: []string{"", " First Name", ""}},
		{Type: token.IDENT, Literal: "Smith", LeadingTrivia: []string{"", "", " Last Name", ""}},
		{Type: token.EOF, LeadingTrivia: []string{"", "", " Final comment"}},
	}, compareLeadingTrivia())
}

func TestEmptySinglelineComment(t *testing.T) {
	assertTokens(t, "//\nhello//\r\nthere//\r!//", []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{"", ""}},
		{Type: token.IDENT, Literal: "there", LeadingTrivia: []string{"", ""}},
		{Type: token.NOT, Literal: "!", LeadingTrivia: []string{"", ""}},
		{Type: token.EOF, Literal: "", LeadingTrivia: []string{""}},
	}, compareLeadingTrivia())
}

func TestLastLineComment(t *testing.T) {
	assertTokens(t, "// last comment", []token.Token{
		{Type: token.EOF, Literal: "", LeadingTrivia: []string{" last comment"}, AfterNewline: false},
	}, compareLeadingTrivia(), compareAfterNewline())
}

func TestScanContinuesAfterNullCharacter(t *testing.T) {
	assertTokens(t, "Hello\x00World", []token.Token{
		{Type: token.IDENT, Literal: "Hello"},
		{Type: token.UNKNOWN, Literal: "\x00"},
		{Type: token.IDENT, Literal: "World"},
		{Type: token.EOF},
	})
}

func TestPunctuators(t *testing.T) {
	assertTokens(t, "; = == ! != < <= > >= () {} /", []token.Token{
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
	assertTokens(t, "  one\ntwo\rthree\tfour \r\n five ", []token.Token{
		{Type: token.IDENT, Literal: "one"},
		{Type: token.IDENT, Literal: "two"},
		{Type: token.IDENT, Literal: "three"},
		{Type: token.IDENT, Literal: "four"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.EOF},
	})
}

func TestReadIdent(t *testing.T) {
	assertTokens(t, " hello  hello123   _hello123 ", []token.Token{
		{Type: token.IDENT, Literal: "hello"},
		{Type: token.IDENT, Literal: "hello123"},
		{Type: token.IDENT, Literal: "_hello123"},
		{Type: token.EOF},
	})
}

func TestReadNumber(t *testing.T) {
	assertTokens(t, "123", []token.Token{
		{Type: token.NUMBER, Literal: "123"},
		{Type: token.EOF},
	})
}

func TestReadString(t *testing.T) {
	t.Run("legal string", func(t *testing.T) {
		assertTokens(t, " 'Hello, World!' \"Hello, World!\"", []token.Token{
			{Type: token.STRING, Literal: "'Hello, World!'"},
			{Type: token.STRING, Literal: "\"Hello, World!\""},
			{Type: token.EOF},
		})
	})
	t.Run("illegal string", func(t *testing.T) {
		inputs := []string{
			"'Hello, World",  // missing '
			"'",              // missing '
			"\"Hello, World", // missing "
			"\"",             // missing "
		}
		assertTokens(t, strings.Join(inputs, "\n"), []token.Token{
			{Type: token.ILLEGAL, Literal: "'Hello, World"},
			{Type: token.ILLEGAL, Literal: "'"},
			{Type: token.ILLEGAL, Literal: "\"Hello, World"},
			{Type: token.ILLEGAL, Literal: "\""},
			{Type: token.EOF},
		})
	})
	t.Run("illegal string with CR in the middle", func(t *testing.T) {
		delimiters := []string{"'", "\""}
		terminators := []string{"\n", "\r", "\r\n"}
		for _, delimiter := range delimiters {
			for _, terminator := range terminators {
				input := fmt.Sprintf("%sHello%sWorld%s", delimiter, terminator, delimiter)
				assertTokens(t, input, []token.Token{
					{Type: token.ILLEGAL, Literal: delimiter + "Hello"},
					{Type: token.IDENT, Literal: "World"},
					{Type: token.ILLEGAL, Literal: delimiter},
					{Type: token.EOF},
				})
			}
		}
	})
}
