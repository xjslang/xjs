package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/token"
)

func TestSkipWhitespaces(t *testing.T) {
	idNames := []string{"lorem", "ipsum", "dolor"}
	l := New(strings.NewReader(fmt.Sprintf("  %s    %s %s   ", idNames[0], idNames[1], idNames[2])))
	for i := range 3 {
		tok := l.NextToken()
		if tok.Literal != idNames[i] {
			t.Errorf("Expected %s, got %s", idNames[i], tok.Literal)
		}
	}
	tok := l.NextToken()
	if tok.Literal != "" {
		t.Errorf("Expected empty string, got %s", tok.Literal)
	}
}

func TestReadIden(t *testing.T) {
	idNames := []string{"hello", "hello123", "_hello123"}
	for _, idName := range idNames {
		l := New(strings.NewReader(idName))
		tok := l.NextToken()
		if tok.Type != token.IDENT {
			t.Errorf("Expected %v, got %v", token.IDENT, tok.Type)
		} else if tok.Literal != idName {
			t.Errorf("Expected %s, got %s", idName, tok.Literal)
		}
	}
}

func TestReadNumber(t *testing.T) {
	number := "123"
	l := New(strings.NewReader(number))
	tok := l.NextToken()
	if tok.Type != token.NUMBER {
		t.Errorf("Expected %v, got %v", token.NUMBER, tok.Type)
	} else if tok.Literal != number {
		t.Errorf("Expected %s, got %s", number, tok.Literal)
	}
}

func TestReadString(t *testing.T) {
	inputs := []string{
		"'Hello, World!'",   // single quote
		"\"Hello, World!\"", // double quote
	}
	for _, input := range inputs {
		l := New(strings.NewReader(input))
		tok := l.NextToken()
		if tok.Type != token.STRING {
			t.Errorf("Expected %v, got %v", token.STRING, tok.Type)
		} else if tok.Literal != input {
			t.Errorf("Expected %s, got %s", input, tok.Literal)
		}
		tok = l.NextToken()
		if tok.Type != token.EOF {
			t.Errorf("Expected %v, got %v", token.EOF, tok.Type)
		}
	}
}

func TestScanContinuesAfterNullCharacter(t *testing.T) {
	l := New(strings.NewReader("hello\x00dolly"))
	expected := []token.Token{
		{Type: token.IDENT, Literal: "hello"},
		{Type: token.UNKNOWN, Literal: "\x00"},
		{Type: token.IDENT, Literal: "dolly"},
	}
	for _, expectedToken := range expected {
		tok := l.NextToken()
		if tok.Type != expectedToken.Type {
			t.Errorf("Expected %v, got %v", expectedToken.Type, tok.Type)
		} else if tok.Literal != expectedToken.Literal {
			t.Errorf("Expected %s, got %s", expectedToken.Literal, tok.Literal)
		}
	}
	tok := l.NextToken()
	if tok.Type != token.EOF {
		t.Errorf("Expected %v, got %v", token.EOF, tok.Type)
	}
}
