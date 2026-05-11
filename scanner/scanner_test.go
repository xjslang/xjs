package scanner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/scanner"
)

func assertLexerTokens(t *testing.T, sc *scanner.Scanner, expectedToks []scanner.Token, opts ...testutil.TokenCompareOption) {
	var toks []scanner.Token
	for {
		tok := sc.NextToken()
		toks = append(toks, tok)
		if tok.Type == scanner.EOF {
			break
		}
	}
	testutil.AssertTokens(t, toks, expectedToks, opts...)
}

func assertInputTokens(t *testing.T, input string, expectedToks []scanner.Token, opts ...testutil.TokenCompareOption) {
	sc := &scanner.Scanner{}
	sc.Init([]byte(input))
	assertLexerTokens(t, sc, expectedToks, opts...)
}

func BenchmarkLexer(b *testing.B) {
	sc := &scanner.Scanner{}
	sc.Init([]byte("lorem ipsum dolor"))
	var tok scanner.Token // prevent dead code elimination
	for b.Loop() {
		for tok = sc.NextToken(); tok.Type != scanner.EOF; tok = sc.NextToken() {
		}
		sc.Reset()
	}
	_ = tok
}

func TestTokenPosition(t *testing.T) {
	input := " aaa   bbb /* block comment*/ ccc\n// comment\rddd\r\ne!\n"
	assertInputTokens(t, input, []scanner.Token{
		{Type: scanner.IDENT, Literal: "aaa", Position: scanner.Position{Line: 0, Column: 1}},
		{Type: scanner.IDENT, Literal: "bbb", Position: scanner.Position{Line: 0, Column: 7}},
		{Type: scanner.IDENT, Literal: "ccc", Position: scanner.Position{Line: 0, Column: 30}},
		{Type: scanner.IDENT, Literal: "ddd", Position: scanner.Position{Line: 2, Column: 0}},
		{Type: scanner.IDENT, Literal: "e", Position: scanner.Position{Line: 3, Column: 0}},
		{Type: scanner.NOT, Literal: "!", Position: scanner.Position{Line: 3, Column: 1}},
		{Type: scanner.EOF, Position: scanner.Position{Line: 4, Column: 0}},
	}, testutil.CompareTokenPosition())
}

func TestReset(t *testing.T) {
	items := []string{"lorem", "ipsum", "dolor"}
	sc := &scanner.Scanner{}
	sc.Init([]byte(strings.Join(items, " ")))
	for range 2 {
		var toks []scanner.Token
		for tok := sc.NextToken(); tok.Type != scanner.EOF; tok = sc.NextToken() {
			toks = append(toks, tok)
		}
		if n := len(toks); n != len(items) {
			t.Fatalf("Expected len(toks)=%d, got %d", len(items), n)
		}
		for i, tok := range toks {
			if tok.Type != scanner.IDENT {
				t.Fatalf("token %d: expected type %v, got %v", i, scanner.IDENT, tok.Type)
			}
			if tok.Literal != items[i] {
				t.Fatalf("token %d: expected %q, got %q", i, items[i], tok.Literal)
			}
		}
		sc.Reset()
	}

	t.Run("without init", func(t *testing.T) {
		sc := &scanner.Scanner{}
		sc.Reset()
		assertLexerTokens(t, sc, []scanner.Token{
			{Type: scanner.EOF},
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
			var expectedToks []scanner.Token
			for _, item := range test.items {
				expectedToks = append(expectedToks, scanner.Token{Type: scanner.STRING, Literal: fmt.Sprintf("'%s'", item)})
			}
			expectedToks = append(expectedToks, scanner.Token{Type: scanner.EOF})
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
			assertInputTokens(t, test.input, []scanner.Token{
				{Type: scanner.IDENT, Literal: "hello"},
				{Type: scanner.IDENT, Literal: "world", AfterNewline: true},
				{Type: scanner.EOF},
			}, testutil.CompareAfterNewline())
		})
	}
}

func TestBlockComments(t *testing.T) {
	input := `/* lorem
ipsum dolor */

hello/* unfinished comment`
	assertInputTokens(t, input, []scanner.Token{
		{Type: scanner.IDENT, Literal: "hello", LeadingTrivia: []string{" lorem\nipsum dolor ", "", ""}},
		{Type: scanner.ILLEGAL, Literal: " unfinished comment"},
		{Type: scanner.EOF},
	}, testutil.CompareLeadingTrivia())
}

func TestLineComments(t *testing.T) {
	input := `
  // First Name
  John
  
  // Last Name
  Smith
	
	// Final comment`
	assertInputTokens(t, input, []scanner.Token{
		{Type: scanner.IDENT, Literal: "John", LeadingTrivia: []string{"", " First Name", ""}},
		{Type: scanner.IDENT, Literal: "Smith", LeadingTrivia: []string{"", "", " Last Name", ""}},
		{Type: scanner.EOF, LeadingTrivia: []string{"", "", " Final comment"}},
	}, testutil.CompareLeadingTrivia())
}

func TestEmptySinglelineComment(t *testing.T) {
	assertInputTokens(t, "//\nhello//\r\nthere//\r!//", []scanner.Token{
		{Type: scanner.IDENT, Literal: "hello", LeadingTrivia: []string{"", ""}},
		{Type: scanner.IDENT, Literal: "there", LeadingTrivia: []string{"", ""}},
		{Type: scanner.NOT, Literal: "!", LeadingTrivia: []string{"", ""}},
		{Type: scanner.EOF, Literal: "", LeadingTrivia: []string{""}},
	}, testutil.CompareLeadingTrivia())
}

func TestLastLineComment(t *testing.T) {
	assertInputTokens(t, "// last comment", []scanner.Token{
		{Type: scanner.EOF, Literal: "", LeadingTrivia: []string{" last comment"}, AfterNewline: false},
	}, testutil.CompareLeadingTrivia(), testutil.CompareAfterNewline())
}

func TestScanContinuesAfterNullCharacter(t *testing.T) {
	assertInputTokens(t, "Hello\x00World", []scanner.Token{
		{Type: scanner.IDENT, Literal: "Hello"},
		{Type: scanner.UNKNOWN, Literal: "\x00"},
		{Type: scanner.IDENT, Literal: "World"},
		{Type: scanner.EOF},
	})
}

func TestPunctuators(t *testing.T) {
	assertInputTokens(t, "; = == ! != < <= > >= () {} + - * / %", []scanner.Token{
		{Type: scanner.SEMICOLON, Literal: ";"},
		{Type: scanner.ASSIGN, Literal: "="},
		{Type: scanner.EQ, Literal: "=="},
		{Type: scanner.NOT, Literal: "!"},
		{Type: scanner.NOT_EQ, Literal: "!="},
		{Type: scanner.LT, Literal: "<"},
		{Type: scanner.LTE, Literal: "<="},
		{Type: scanner.GT, Literal: ">"},
		{Type: scanner.GTE, Literal: ">="},
		{Type: scanner.LPAREN, Literal: "("},
		{Type: scanner.RPAREN, Literal: ")"},
		{Type: scanner.LBRACE, Literal: "{"},
		{Type: scanner.RBRACE, Literal: "}"},
		{Type: scanner.PLUS, Literal: "+"},
		{Type: scanner.MINUS, Literal: "-"},
		{Type: scanner.MULTIPLY, Literal: "*"},
		{Type: scanner.DIVIDE, Literal: "/"},
		{Type: scanner.MODULO, Literal: "%"},
		{Type: scanner.EOF},
	})
}

func TestSkipWhitespaces(t *testing.T) {
	assertInputTokens(t, "  one\ntwo\rthree\tfour \r\n five ", []scanner.Token{
		{Type: scanner.IDENT, Literal: "one"},
		{Type: scanner.IDENT, Literal: "two"},
		{Type: scanner.IDENT, Literal: "three"},
		{Type: scanner.IDENT, Literal: "four"},
		{Type: scanner.IDENT, Literal: "five"},
		{Type: scanner.EOF},
	})
}

func TestReadIdent(t *testing.T) {
	assertInputTokens(t, " hello  hello123   _hello123 ", []scanner.Token{
		{Type: scanner.IDENT, Literal: "hello"},
		{Type: scanner.IDENT, Literal: "hello123"},
		{Type: scanner.IDENT, Literal: "_hello123"},
		{Type: scanner.EOF},
	})
}

func TestReadNumber(t *testing.T) {
	assertInputTokens(t, "123", []scanner.Token{
		{Type: scanner.NUMBER, Literal: "123"},
		{Type: scanner.EOF},
	})
}

func TestReadBoolean(t *testing.T) {
	assertInputTokens(t, "true false", []scanner.Token{
		{Type: scanner.BOOLEAN, Literal: "true"},
		{Type: scanner.BOOLEAN, Literal: "false"},
		{Type: scanner.EOF},
	})
}

func TestReadString(t *testing.T) {
	t.Run("legal string", func(t *testing.T) {
		assertInputTokens(t, " 'Hello, World!' \"Hello, World!\"", []scanner.Token{
			{Type: scanner.STRING, Literal: "'Hello, World!'"},
			{Type: scanner.STRING, Literal: "\"Hello, World!\""},
			{Type: scanner.EOF},
		})
	})
	t.Run("illegal string", func(t *testing.T) {
		inputs := []string{
			"'Hello, World",  // missing '
			"'",              // missing '
			"\"Hello, World", // missing "
			"\"",             // missing "
		}
		assertInputTokens(t, strings.Join(inputs, "\n"), []scanner.Token{
			{Type: scanner.ILLEGAL, Literal: "'Hello, World"},
			{Type: scanner.ILLEGAL, Literal: "'"},
			{Type: scanner.ILLEGAL, Literal: "\"Hello, World"},
			{Type: scanner.ILLEGAL, Literal: "\""},
			{Type: scanner.EOF},
		})
	})
	t.Run("illegal string with CR in the middle", func(t *testing.T) {
		delimiters := []string{"'", "\""}
		terminators := []string{"\n", "\r", "\r\n"}
		for _, delimiter := range delimiters {
			for _, terminator := range terminators {
				input := fmt.Sprintf("%sHello%sWorld%s", delimiter, terminator, delimiter)
				assertInputTokens(t, input, []scanner.Token{
					{Type: scanner.ILLEGAL, Literal: delimiter + "Hello"},
					{Type: scanner.IDENT, Literal: "World"},
					{Type: scanner.ILLEGAL, Literal: delimiter},
					{Type: scanner.EOF},
				})
			}
		}
	})
}

func TestKeywords(t *testing.T) {
	input := "let function"
	assertInputTokens(t, input, []scanner.Token{
		{Type: scanner.LET, Literal: "let"},
		{Type: scanner.FUNCTION, Literal: "function"},
		{Type: scanner.EOF},
	})
}
