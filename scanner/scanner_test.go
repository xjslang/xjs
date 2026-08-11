package scanner_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type tokenCompareConfig struct {
	afterNewline  bool
	leadingTrivia bool
	tokenPosition bool
}

type TokenCompareOption func(cfg *tokenCompareConfig)

func CompareAfterNewline() TokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.afterNewline = true
	}
}

func CompareLeadingTrivia() TokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.leadingTrivia = true
	}
}

func CompareTokenPosition() TokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.tokenPosition = true
	}
}

func AssertTokens(t *testing.T, toks, expectedToks []token.Token, opts ...TokenCompareOption) {
	t.Helper()
	cfg := &tokenCompareConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if len(toks) != len(expectedToks) {
		t.Fatalf("Expect len(toks) = %d, got %d", len(toks), len(expectedToks))
	}
	for i, expectedTok := range expectedToks {
		tok := toks[i]
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
				for j, expectedTrivia := range expectedTok.LeadingTrivia {
					trivia := tok.LeadingTrivia[j]
					if trivia.Type != expectedTrivia.Type {
						t.Errorf("token %d: expected trivia type to be %v, got %v", i, expectedTrivia.Type, trivia.Type)
					} else if trivia.Literal != expectedTrivia.Literal {
						t.Errorf("token %d: expected trivia to be %q, got %q", i, expectedTrivia.Literal, trivia.Literal)
					}
				}
			}
		case cfg.tokenPosition && (tok.Line != expectedTok.Line || tok.Column != expectedTok.Column):
			t.Errorf("token %d: expected position to be (%d, %d), got (%d, %d)", i, expectedTok.Line, expectedTok.Column, tok.Line, tok.Column)
		}
	}
}

func assertLexerTokens(t *testing.T, sc *scanner.Scanner, expectedToks []token.Token, opts ...TokenCompareOption) {
	t.Helper()
	var toks []token.Token
	for {
		tok := sc.NextToken()
		toks = append(toks, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	AssertTokens(t, toks, expectedToks, opts...)
}

func assertInputTokens(t *testing.T, input string, expectedToks []token.Token, opts ...TokenCompareOption) {
	t.Helper()
	sb := scanner.Builder{}
	s := sb.Build([]byte(input))
	assertLexerTokens(t, s, expectedToks, opts...)
}

func ExampleBuilder_Build() {
	hashTyp := token.RegisterType("HASH", "#")
	caretType := token.RegisterType("CARET", "^")
	sb := scanner.Builder{}
	s := sb.
		UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
			if sc.CurrentChar() == '#' {
				sc.AdvanceChar()
				return token.Token{Type: hashTyp, Literal: "#"}, nil
			}
			return next() // Delegate to the "next" middleware
		}).
		UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
			if sc.CurrentChar() == '^' {
				sc.AdvanceChar()
				return token.Token{Type: caretType, Literal: "^"}, nil
			}
			return next() // Delegate to the "next" middleware
		}).
		Build([]byte("#some ^input"))

	// Now you can use the scanner
	for tok := s.NextToken(); tok.Type != token.EOF; tok = s.NextToken() {
		fmt.Printf(
			"{Type: %s, Literal: %s, Position: %v}\n",
			tok.Type.Name(), tok.Literal, tok.Position)
	}
	// Output:
	// {Type: HASH, Literal: #, Position: {0 0 0}}
	// {Type: IDENT, Literal: some, Position: {0 1 1}}
	// {Type: CARET, Literal: ^, Position: {0 6 6}}
	// {Type: IDENT, Literal: input, Position: {0 7 8}}
}

func BenchmarkUseScanner(b *testing.B) {
	dat, err := os.ReadFile("testdata/sample.js")
	if err != nil {
		b.Fatal(err)
	}

	sb := scanner.Builder{}
	for range 10 {
		sb.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
			if tok, err = next(); err != nil {
				return
			}
			return
		})
	}

	s := sb.Build(dat)
	b.ResetTimer()
	for b.Loop() {
		s.Reset()
		for tok := s.NextToken(); tok.Type != token.EOF; tok = s.NextToken() {
		}
	}
}

func BenchmarkLexer(b *testing.B) {
	sb := scanner.Builder{}
	sc := sb.Build([]byte("lorem ipsum dolor"))
	var tok token.Token // prevent dead code elimination
	for b.Loop() {
		for tok = sc.NextToken(); tok.Type != token.EOF; tok = sc.NextToken() {
		}
		sc.Reset()
	}
	_ = tok
}

func TestLookahead(t *testing.T) {
	input := "a b c d e f"
	sb := scanner.Builder{}
	sc := sb.Build([]byte(input))
	var toks []token.Token
	for range 2 {
		toks = append(toks, sc.NextToken())
	}
	AssertTokens(t, toks, []token.Token{
		{Type: token.IDENT, Literal: "a"},
		{Type: token.IDENT, Literal: "b"},
	})
	fork := sc.Fork()
	toks = toks[:0]
	for range 2 {
		toks = append(toks, fork.NextToken())
	}
	AssertTokens(t, toks, []token.Token{
		{Type: token.IDENT, Literal: "c"},
		{Type: token.IDENT, Literal: "d"},
	})
	toks = toks[:0]
	toks = append(toks, sc.NextToken())
	AssertTokens(t, toks, []token.Token{
		{Type: token.IDENT, Literal: "c"},
	})
	sc.Apply(fork)
	toks = toks[:0]
	for range 2 {
		toks = append(toks, sc.NextToken())
	}
	AssertTokens(t, toks, []token.Token{
		{Type: token.IDENT, Literal: "e"},
		{Type: token.IDENT, Literal: "f"},
	})
}

func TestIdentifier(t *testing.T) {
	input := "abc _abc $abc Abc0123 _$AbC0$__ 0abc"
	assertInputTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "abc"},
		{Type: token.IDENT, Literal: "_abc"},
		{Type: token.IDENT, Literal: "$abc"},
		{Type: token.IDENT, Literal: "Abc0123"},
		{Type: token.IDENT, Literal: "_$AbC0$__"},
		{Type: token.ILLEGAL, Literal: "0abc"},
		{Type: token.EOF},
	})
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
	}, CompareTokenPosition())
}

func TestReset(t *testing.T) {
	items := []string{"lorem", "ipsum", "dolor"}
	sb := scanner.Builder{}
	sc := sb.Build([]byte(strings.Join(items, " ")))
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
			}, CompareAfterNewline())
		})
	}
}

func TestBlockComments(t *testing.T) {
	input := "/* lorem\nipsum dolor */\n\rhello\r\n/* unfinished comment"
	assertInputTokens(t, input, []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []token.Token{
			{Type: token.BLOCK_COMMENT, Literal: "/* lorem\nipsum dolor */"},
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.NEWLINE, Literal: "\r"},
		}},
		{Type: token.ILLEGAL, Literal: "/* unfinished comment", LeadingTrivia: []token.Token{
			{Type: token.NEWLINE, Literal: "\r\n"},
		}},
		{Type: token.EOF},
	}, CompareLeadingTrivia())
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
			{Type: token.LINE_COMMENT, Literal: "// First Name\n"},
		}},
		{Type: token.IDENT, Literal: "Smith", LeadingTrivia: []token.Token{
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.LINE_COMMENT, Literal: "// Last Name\n"},
		}},
		{Type: token.EOF, LeadingTrivia: []token.Token{
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.NEWLINE, Literal: "\n"},
			{Type: token.LINE_COMMENT, Literal: "// Final comment"},
		}},
	}, CompareLeadingTrivia())
}

func TestEmptySinglelineComment(t *testing.T) {
	assertInputTokens(t, "//\nhello//\n\npeople//\r\nthere//\r!//", []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "//\n"},
		}},
		{Type: token.IDENT, Literal: "people", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "//\n"},
			{Type: token.NEWLINE, Literal: "\n"},
		}},
		{Type: token.IDENT, Literal: "there", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "//\r\n"},
		}},
		{Type: token.NOT, Literal: "!", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "//\r"},
		}},
		{Type: token.EOF, Literal: "", LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "//"},
		}},
	}, CompareLeadingTrivia())
}

func TestLastLineComment(t *testing.T) {
	assertInputTokens(t, "// last comment", []token.Token{
		{Type: token.EOF, Literal: "", AfterNewline: false, LeadingTrivia: []token.Token{
			{Type: token.LINE_COMMENT, Literal: "// last comment"},
		}},
	}, CompareLeadingTrivia(), CompareAfterNewline())
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
	assertInputTokens(t, "123 0.5 0e2 0123 0x10 0o7", []token.Token{
		{Type: token.NUMBER, Literal: "123"},
		{Type: token.NUMBER, Literal: "0.5"},
		{Type: token.NUMBER, Literal: "0e2"},
		{Type: token.NUMBER, Literal: "0123"},
		{Type: token.NUMBER, Literal: "0x10"},
		{Type: token.NUMBER, Literal: "0o7"},
		{Type: token.EOF},
	})
}

func TestReadString(t *testing.T) {
	t.Run("legal string", func(t *testing.T) {
		assertInputTokens(t, " 'Hello, World!' \"Hello, World!\" `Hello,\nWorld!`", []token.Token{
			{Type: token.STRING, Literal: "'Hello, World!'"},
			{Type: token.STRING, Literal: "\"Hello, World!\""},
			{Type: token.STRING, Literal: "`Hello,\nWorld!`"},
			{Type: token.EOF},
		})
	})
	t.Run("illegal string", func(t *testing.T) {
		inputs := []string{
			"'Hello, World",  // missing '
			"'",              // missing '
			"\"Hello, World", // missing "
			"\"",             // missing "
			"`Hello,\nWorld", // missing `
		}
		assertInputTokens(t, strings.Join(inputs, "\n"), []token.Token{
			{Type: token.ILLEGAL, Literal: "'Hello, World"},
			{Type: token.ILLEGAL, Literal: "'"},
			{Type: token.ILLEGAL, Literal: "\"Hello, World"},
			{Type: token.ILLEGAL, Literal: "\""},
			{Type: token.ILLEGAL, Literal: "`Hello,\nWorld"},
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

func TestUseScanner(t *testing.T) {
	powType := token.RegisterType("POW", "**")
	sb := scanner.Builder{}
	sc := sb.
		UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
			if sc.CurrentChar() == '*' && sc.PeekChar() == '*' {
				// consume **
				sc.AdvanceChar()
				sc.AdvanceChar()
				return token.Token{Type: powType, Literal: powType.Literal()}, nil
			}
			return next()
		}).
		Build([]byte("5 ** 2"))
	assertLexerTokens(t, sc, []token.Token{
		{Type: token.NUMBER, Literal: "5"},
		{Type: powType, Literal: "**"},
		{Type: token.NUMBER, Literal: "2"},
		{Type: token.EOF},
	})
}

func TestScanNumber(t *testing.T) {
	t.Run("invalid formats", func(t *testing.T) {
		inputs := []string{
			"123e", "123e+", "123e-", "1e", // invalid float numbers
			"0x", // invalid hex number
			"0o", // invalid octal number
		}
		for _, input := range inputs {
			sb := scanner.Builder{}
			sc := sb.Build([]byte(input))
			_, err := scanner.ScanNumber(sc)
			assert.Error(t, err)
		}
	})
}

func TestSpecialRunes(t *testing.T) {
	inputs := []string{"T\u200c", "A\u200d"}
	for _, input := range inputs {
		sb := scanner.Builder{}
		sc := sb.Build([]byte(input))
		tok := sc.NextToken()
		require.Equal(t, token.IDENT, tok.Type)
		require.Equal(t, input, tok.Literal)
		require.Equal(t, token.EOF, sc.NextToken().Type)
	}
}
