package scanner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func assertLexerTokens(t *testing.T, sc *scanner.Scanner, expectedToks []token.Token, opts ...testutil.TokenCompareOption) {
	t.Helper()
	var toks []token.Token
	for {
		tok := sc.NextToken()
		toks = append(toks, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	testutil.AssertTokens(t, toks, expectedToks, opts...)
}

func assertInputTokens(t *testing.T, input string, expectedToks []token.Token, opts ...testutil.TokenCompareOption) {
	t.Helper()
	sc := &scanner.Scanner{}
	sc.Init([]byte(input))
	assertLexerTokens(t, sc, expectedToks, opts...)
}

func ExampleScanner_Init() {
	hashTyp := token.RegisterType("hash")
	caretType := token.RegisterType("caret")
	s := &scanner.Scanner{}

	// Declare "middlewares" BEFORE calling Init
	s.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
		if sc.CurrentChar() == '#' {
			sc.AdvanceChar()
			return token.Token{Type: hashTyp, Literal: "#"}, nil
		}
		return next() // Delegate to the "next" middleware
	})
	s.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
		if sc.CurrentChar() == '^' {
			sc.AdvanceChar()
			return token.Token{Type: caretType, Literal: "^"}, nil
		}
		return next() // Delegate to the "next" middleware
	})
	s.Init([]byte("#some ^input"))

	// Now you can use the scanner
	for tok := s.NextToken(); tok.Type != token.EOF; tok = s.NextToken() {
		fmt.Printf(
			"{Type: %s, Literal: %s, Position: %v}\n",
			tok.Type.String(),
			tok.Literal, tok.Position)
	}
	// Output:
	// {Type: hash, Literal: #, Position: {0 0}}
	// {Type: identifier, Literal: some, Position: {0 1}}
	// {Type: caret, Literal: ^, Position: {0 6}}
	// {Type: identifier, Literal: input, Position: {0 7}}
}

func BenchmarkLexer(b *testing.B) {
	sc := &scanner.Scanner{}
	sc.Init([]byte("lorem ipsum dolor"))
	var tok token.Token // prevent dead code elimination
	for b.Loop() {
		for tok = sc.NextToken(); tok.Type != token.EOF; tok = sc.NextToken() {
		}
		sc.Reset()
	}
	_ = tok
}

func TestTokenPosition(t *testing.T) {
	input := " aaa   bbb /* block comment*/ ccc\n// comment\rddd\r\ne!\n"
	assertInputTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "aaa", Position: token.Position{Line: 0, Column: 1}},
		{Type: token.IDENT, Literal: "bbb", Position: token.Position{Line: 0, Column: 7}},
		{Type: token.IDENT, Literal: "ccc", Position: token.Position{Line: 0, Column: 30}},
		{Type: token.IDENT, Literal: "ddd", Position: token.Position{Line: 2, Column: 0}},
		{Type: token.IDENT, Literal: "e", Position: token.Position{Line: 3, Column: 0}},
		{Type: token.NOT, Literal: "!", Position: token.Position{Line: 3, Column: 1}},
		{Type: token.EOF, Position: token.Position{Line: 4, Column: 0}},
	}, testutil.CompareTokenPosition())
}

func TestReset(t *testing.T) {
	items := []string{"lorem", "ipsum", "dolor"}
	sc := &scanner.Scanner{}
	sc.Init([]byte(strings.Join(items, " ")))
	for range 2 {
		var toks []token.Token
		for tok := sc.NextToken(); tok.Type != token.EOF; tok = sc.NextToken() {
			toks = append(toks, tok)
		}
		if n := len(toks); n != len(items) {
			t.Fatalf("Expected len(toks)=%d, got %d", len(items), n)
		}
		for i, tok := range toks {
			if tok.Type != token.IDENT {
				t.Fatalf("token %d: expected type %v, got %v", i, token.IDENT, tok.Type)
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
		assertLexerTokens(t, sc, []token.Token{
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
	input := "/* lorem\nipsum dolor */\n\rhello\r\n/* unfinished comment"
	assertInputTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []token.Token{
			{Type: token.BLOCK_COMMENT, Literal: " lorem\nipsum dolor "},
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.NEWLINE, Literal: "\r"},
		}},
		{Type: token.ILLEGAL, Literal: " unfinished comment", LeadingTrivia: []token.Token{
			{Type: token.NEWLINE, Literal: "\r\n"},
		}},
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
		{Type: token.IDENT, Literal: "John", LeadingTrivia: []token.Token{
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.LINE_COMMENT, Literal: " First Name\n"},
		}},
		{Type: token.IDENT, Literal: "Smith", LeadingTrivia: []token.Token{
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.LINE_COMMENT, Literal: " Last Name\n"},
		}},
		{Type: token.EOF, LeadingTrivia: []token.Token{
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.LINE_COMMENT, Literal: " Final comment"},
		}},
	}, testutil.CompareLeadingTrivia())
}

func TestEmptySinglelineComment(t *testing.T) {
	assertInputTokens(t, "//\nhello//\n\npeople//\r\nthere//\r!//", []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "\n"},
		}},
		{Type: token.IDENT, Literal: "people", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "\n"},
			{Type: token.NEWLINE, Literal: "\n"},
		}},
		{Type: token.IDENT, Literal: "there", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "\r\n"},
		}},
		{Type: token.NOT, Literal: "!", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "\r"},
		}},
		{Type: token.EOF, Literal: "", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT},
		}},
	}, testutil.CompareLeadingTrivia())
}

func TestLastLineComment(t *testing.T) {
	assertInputTokens(t, "// last comment", []token.Token{
		{Type: token.EOF, Literal: "", AfterNewline: false, LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: " last comment"},
		}},
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
	assertInputTokens(t, "; = == ! != < <= > >= () {} + ++ - -- * / % && || | &", []token.Token{
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
		{Type: token.PLUS, Literal: "+"},
		{Type: token.INCREMENT, Literal: "++"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.DECREMENT, Literal: "--"},
		{Type: token.MULTIPLY, Literal: "*"},
		{Type: token.DIVIDE, Literal: "/"},
		{Type: token.MODULO, Literal: "%"},
		{Type: token.AND, Literal: "&&"},
		{Type: token.OR, Literal: "||"},
		{Type: token.UNKNOWN, Literal: "|"},
		{Type: token.UNKNOWN, Literal: "&"},
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
