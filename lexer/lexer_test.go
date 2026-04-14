package lexer

import (
	"strings"
	"testing"

	"github.com/xjslang/xjs/token"
)

func TestAfterNewline(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"newline before block comment", "hello\n/* block comment */world"},
		{"block comment with newline", "hello/* block\ncomment */world"},
		{"newline", "hello\nworld"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedToks := []token.Token{
				{Type: token.IDENT, Literal: "hello"},
				{Type: token.IDENT, Literal: "world", AfterNewline: true},
			}
			l := New(strings.NewReader(test.input))
			for i, expectedTok := range expectedToks {
				tok := l.NextToken()
				switch {
				case tok.Type != expectedTok.Type:
					t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
				case tok.Literal != expectedTok.Literal:
					t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
				case tok.AfterNewline != expectedTok.AfterNewline:
					t.Errorf("token %d: expected AfterNewline to be %t, got %t", i, expectedTok.AfterNewline, tok.AfterNewline)
				}
			}
		})
	}
}

func TestEmptySinglelineComment(t *testing.T) {
	expectedToks := []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{""}},
		{Type: token.IDENT, Literal: "there", LeadingTrivia: []string{""}},
		{Type: token.EOF, LeadingTrivia: []string{"\rObi-Wan Kenobi"}},
	}
	l := New(strings.NewReader("//\nhello//\r\nthere//\rObi-Wan Kenobi"))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		case len(tok.LeadingTrivia) != len(expectedTok.LeadingTrivia):
			t.Errorf("token %d: expected %d leading trivia lines, got %d", i, len(expectedTok.LeadingTrivia), len(tok.LeadingTrivia))
		default:
			for j, line := range expectedTok.LeadingTrivia {
				if tok.LeadingTrivia[j] != line {
					t.Errorf("token %d: expected %q leading trivia line, got %q", i, line, tok.LeadingTrivia[j])
				}
			}
		}
	}
}

func TestMultilineComments(t *testing.T) {
	expectedToks := []token.Token{
		{Type: token.IDENT, Literal: "hello", LeadingTrivia: []string{" lorem\nipsum dolor ", "", ""}},
		{Type: token.ILLEGAL, Literal: " unfinished comment"},
		{Type: token.EOF},
	}
	l := New(strings.NewReader(`/* lorem
ipsum dolor */

hello/* unfinished comment`))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		case len(tok.LeadingTrivia) != len(expectedTok.LeadingTrivia):
			t.Errorf("token %d: expected %d leading trivia lines, got %d", i, len(expectedTok.LeadingTrivia), len(tok.LeadingTrivia))
		default:
			for j, line := range expectedTok.LeadingTrivia {
				if tok.LeadingTrivia[j] != line {
					t.Errorf("token %d: expected %q leading trivia line, got %q", i, line, tok.LeadingTrivia[j])
				}
			}
		}
	}
}

func TestSinglelineComments(t *testing.T) {
	expectedToks := []token.Token{
		{Type: token.IDENT, Literal: "John", LeadingTrivia: []string{"", " First Name"}},
		{Type: token.IDENT, Literal: "Smith", LeadingTrivia: []string{"", "", " Last Name"}},
		{Type: token.EOF, LeadingTrivia: []string{"", "", " Final comment"}},
	}
	l := New(strings.NewReader(`
  // First Name
  John
  
  // Last Name
  Smith
	
	// Final comment`))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		case len(tok.LeadingTrivia) != len(expectedTok.LeadingTrivia):
			t.Errorf("token %d: expected %d leading trivia lines, got %d", i, len(expectedTok.LeadingTrivia), len(tok.LeadingTrivia))
		default:
			for j, line := range expectedTok.LeadingTrivia {
				if tok.LeadingTrivia[j] != line {
					t.Errorf("token %d: expected %q leading trivia line, got %q", i, line, tok.LeadingTrivia[j])
				}
			}
		}
	}
}

func TestScanContinuesAfterNullCharacter(t *testing.T) {
	expectedToks := []token.Token{
		{Type: token.IDENT, Literal: "Hello"},
		{Type: token.UNKNOWN, Literal: "\x00"},
		{Type: token.IDENT, Literal: "World"},
		{Type: token.EOF},
	}
	l := New(strings.NewReader("Hello\x00World"))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		}
	}
}

func TestPunctuators(t *testing.T) {
	expectedToks := []token.Token{
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
	}
	l := New(strings.NewReader("; = == ! != < <= > >= () {} /"))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		}
	}
}

func TestSkipWhitespaces(t *testing.T) {
	expectedToks := []token.Token{
		{Type: token.IDENT, Literal: "one"},
		{Type: token.IDENT, Literal: "two"},
		{Type: token.IDENT, Literal: "three"},
		{Type: token.IDENT, Literal: "four"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.EOF},
	}
	l := New(strings.NewReader("  one\ntwo\rthree\tfour \r\n five "))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		}
	}
}

func TestReadIdent(t *testing.T) {
	expectedToks := []token.Token{
		{Type: token.IDENT, Literal: "hello"},
		{Type: token.IDENT, Literal: "hello123"},
		{Type: token.IDENT, Literal: "_hello123"},
		{Type: token.EOF},
	}
	l := New(strings.NewReader(" hello  hello123   _hello123 "))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		}
	}
}

func TestReadNumber(t *testing.T) {
	expectedToks := []token.Token{
		{Type: token.NUMBER, Literal: "123"},
		{Type: token.EOF},
	}
	l := New(strings.NewReader("123"))
	for i, expectedTok := range expectedToks {
		tok := l.NextToken()
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		}
	}
}

func TestReadString(t *testing.T) {
	t.Run("legal string", func(t *testing.T) {
		expectedToks := []token.Token{
			{Type: token.STRING, Literal: "'Hello, World!'"},
			{Type: token.STRING, Literal: "\"Hello, World!\""},
			{Type: token.EOF},
		}
		l := New(strings.NewReader(" 'Hello, World!' \"Hello, World!\""))
		for i, expectedTok := range expectedToks {
			tok := l.NextToken()
			switch {
			case tok.Type != expectedTok.Type:
				t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
			case tok.Literal != expectedTok.Literal:
				t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
			}
		}
	})
	t.Run("illegal string", func(t *testing.T) {
		inputs := []string{
			"'Hello, World",  // missing '
			"'",              // missing '
			"\"Hello, World", // missing "
			"\"",             // missing "
		}
		expectedToks := []token.Token{
			{Type: token.ILLEGAL, Literal: "'Hello, World"},
			{Type: token.ILLEGAL, Literal: "'"},
			{Type: token.ILLEGAL, Literal: "\"Hello, World"},
			{Type: token.ILLEGAL, Literal: "\""},
			{Type: token.EOF},
		}
		l := New(strings.NewReader(strings.Join(inputs, "\n")))
		for i, expectedTok := range expectedToks {
			tok := l.NextToken()
			switch {
			case tok.Type != expectedTok.Type:
				t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
			case tok.Literal != expectedTok.Literal:
				t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
			}
		}
	})
}
