package lexer

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/source"
	"github.com/xjslang/xjs/token"
)

func assertLexerTokens(t *testing.T, l *Lexer, expectedToks []token.Token, opts ...testutil.TokenCompareOption) {
	var toks []token.Token
	for {
		tok := l.NextToken()
		toks = append(toks, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	testutil.AssertTokens(t, toks, expectedToks, opts...)
}

func assertInputTokens(t *testing.T, input string, expectedToks []token.Token, opts ...testutil.TokenCompareOption) {
	l := &Lexer{}
	l.Init([]byte(input))
	assertLexerTokens(t, l, expectedToks, opts...)
}

func BenchmarkLexer(b *testing.B) {
	l := &Lexer{}
	l.Init([]byte("lorem ipsum dolor"))
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
	assertInputTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "aaa", Position: source.Position{Line: 0, Column: 1}},
		{Type: token.IDENT, Literal: "bbb", Position: source.Position{Line: 0, Column: 7}},
		{Type: token.IDENT, Literal: "ccc", Position: source.Position{Line: 0, Column: 30}},
		{Type: token.IDENT, Literal: "ddd", Position: source.Position{Line: 2, Column: 0}},
		{Type: token.IDENT, Literal: "e", Position: source.Position{Line: 3, Column: 0}},
		{Type: token.NOT, Literal: "!", Position: source.Position{Line: 3, Column: 1}},
		{Type: token.EOF, Position: source.Position{Line: 4, Column: 0}},
	}, testutil.CompareTokenPosition())
}

func TestAdvanceChar(t *testing.T) {
	t.Run("column stops at last position", func(t *testing.T) {
		input := "hello"
		l := &Lexer{}
		l.Init([]byte(input))
		for range len(input) {
			l.AdvanceChar()
		}
		if l.column != len(input)-1 {
			t.Errorf("Expected column to be %d, got %d", len(input)-1, l.column)
		}
	})
	t.Run("token column is never negative", func(t *testing.T) {
		assertInputTokens(t, "hello\n", []token.Token{
			{Type: token.IDENT, Literal: "hello", Position: source.Position{Line: 0, Column: 0}},
			{Type: token.EOF, Literal: "", Position: source.Position{Line: 1, Column: 0}},
		}, testutil.CompareTokenPosition())
	})
	t.Run("invalid byte at end is not EOF", func(t *testing.T) {
		assertInputTokens(t, "hello\xff", []token.Token{
			{Type: token.IDENT, Literal: "hello"},
			{Type: token.ILLEGAL, Literal: string(utf8.RuneError)},
			{Type: token.EOF},
		})
	})
}

func TestReset(t *testing.T) {
	items := []string{"lorem", "ipsum", "dolor"}
	l := &Lexer{}
	l.Init([]byte(strings.Join(items, " ")))
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

	t.Run("without init", func(t *testing.T) {
		l := &Lexer{}
		l.Reset()
		assertLexerTokens(t, l, []token.Token{
			{Type: token.EOF},
		})
	})
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
			assertInputTokens(t, item, expectedToks)
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
			assertInputTokens(t, test.input, []token.Token{
				{Type: token.IDENT, Literal: "hello"},
				{Type: token.IDENT, Literal: "world", AfterNewline: true},
				{Type: token.EOF},
			}, testutil.CompareAfterNewline())
		})
	}
}

func TestBlockComments(t *testing.T) {
	input := `/* lorem
ipsum dolor */

hello/* unfinished comment`
	assertInputTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{" lorem\nipsum dolor ", "", ""}},
		{Type: token.ILLEGAL, Literal: " unfinished comment"},
		{Type: token.EOF},
	}, testutil.CompareLeadingTrivia())
}

func TestLineComments(t *testing.T) {
	input := `
  // First Name
  John
  
  // Last Name
  Smith
	
	// Final comment`
	assertInputTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "John", LeadingTrivia: []string{"", " First Name", ""}},
		{Type: token.IDENT, Literal: "Smith", LeadingTrivia: []string{"", "", " Last Name", ""}},
		{Type: token.EOF, LeadingTrivia: []string{"", "", " Final comment"}},
	}, testutil.CompareLeadingTrivia())
}

func TestEmptySinglelineComment(t *testing.T) {
	assertInputTokens(t, "//\nhello//\r\nthere//\r!//", []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{"", ""}},
		{Type: token.IDENT, Literal: "there", LeadingTrivia: []string{"", ""}},
		{Type: token.NOT, Literal: "!", LeadingTrivia: []string{"", ""}},
		{Type: token.EOF, Literal: "", LeadingTrivia: []string{""}},
	}, testutil.CompareLeadingTrivia())
}

func TestLastLineComment(t *testing.T) {
	assertInputTokens(t, "// last comment", []token.Token{
		{Type: token.EOF, Literal: "", LeadingTrivia: []string{" last comment"}, AfterNewline: false},
	}, testutil.CompareLeadingTrivia(), testutil.CompareAfterNewline())
}

func TestScanContinuesAfterNullCharacter(t *testing.T) {
	assertInputTokens(t, "Hello\x00World", []token.Token{
		{Type: token.IDENT, Literal: "Hello"},
		{Type: token.UNKNOWN, Literal: "\x00"},
		{Type: token.IDENT, Literal: "World"},
		{Type: token.EOF},
	})
}

func TestPunctuators(t *testing.T) {
	assertInputTokens(t, "; = == ! != < <= > >= () {} /", []token.Token{
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
	assertInputTokens(t, "  one\ntwo\rthree\tfour \r\n five ", []token.Token{
		{Type: token.IDENT, Literal: "one"},
		{Type: token.IDENT, Literal: "two"},
		{Type: token.IDENT, Literal: "three"},
		{Type: token.IDENT, Literal: "four"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.EOF},
	})
}

func TestReadIdent(t *testing.T) {
	assertInputTokens(t, " hello  hello123   _hello123 ", []token.Token{
		{Type: token.IDENT, Literal: "hello"},
		{Type: token.IDENT, Literal: "hello123"},
		{Type: token.IDENT, Literal: "_hello123"},
		{Type: token.EOF},
	})
}

func TestReadNumber(t *testing.T) {
	assertInputTokens(t, "123", []token.Token{
		{Type: token.NUMBER, Literal: "123"},
		{Type: token.EOF},
	})
}

func TestReadBoolean(t *testing.T) {
	assertInputTokens(t, "true false", []token.Token{
		{Type: token.BOOLEAN, Literal: "true"},
		{Type: token.BOOLEAN, Literal: "false"},
		{Type: token.EOF},
	})
}

func TestReadString(t *testing.T) {
	t.Run("legal string", func(t *testing.T) {
		assertInputTokens(t, " 'Hello, World!' \"Hello, World!\"", []token.Token{
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
		assertInputTokens(t, strings.Join(inputs, "\n"), []token.Token{
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
				assertInputTokens(t, input, []token.Token{
					{Type: token.ILLEGAL, Literal: delimiter + "Hello"},
					{Type: token.IDENT, Literal: "World"},
					{Type: token.ILLEGAL, Literal: delimiter},
					{Type: token.EOF},
				})
			}
		}
	})
}
