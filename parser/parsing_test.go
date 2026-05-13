package parser_test

import (
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
)

func createParser(input string) *parser.Parser {
	l := &scanner.Scanner{}
	l.Init([]byte(input))
	p := &parser.Parser{}
	p.Init(l)
	return p
}

func TestMalformedExpression(t *testing.T) {
	t.Run("let", func(t *testing.T) {
		input := "x = 100"
		p := createParser(input)
		_, err := parser.ParseLet(p)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		expected := "Expected let"
		if got := err.Error(); got != expected {
			t.Fatalf("Expected error to be %q, got %q", expected, got)
		}
	})
	t.Run("block", func(t *testing.T) {
		tests := []struct {
			input       string
			expectedErr string
		}{
			{"let x = 100 }", "Expected {"},
			{"{ let x = 100", "Expected }"},
		}
		for i, test := range tests {
			p := createParser(test.input)
			_, err := parser.ParseBlock(p)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			if got := err.Error(); got != test.expectedErr {
				t.Fatalf("%d: Expected error to be %q, got %q", i, test.expectedErr, got)
			}
		}
	})
	t.Run("function", func(t *testing.T) {
		input := "() {}"
		p := createParser(input)
		_, err := parser.ParseFunction(p)
		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
		expected := "Expected function"
		if got := err.Error(); got != expected {
			t.Fatalf("Expected error to be %q, got %q", expected, got)
		}
	})
	t.Run("grouped expression", func(t *testing.T) {
		tests := []struct {
			input       string
			expectedErr string
		}{
			{"1 + 2)", "Expected ("},
			{"(1 + 2", "Expected )"},
		}
		for i, test := range tests {
			p := createParser(test.input)
			_, err := parser.ParseGroupedExpression(p)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			if got := err.Error(); got != test.expectedErr {
				t.Fatalf("%d: Expected error to be %q, got %q", i, test.expectedErr, got)
			}
		}
	})
}

func TestKeysAreSaved(t *testing.T) {
	t.Run("let", func(t *testing.T) {
		input := `// comment
		let area
		= 200 /*c*/;`
		p := createParser(input)
		result, err := parser.ParseLet(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]scanner.Token{result.LetToken, result.AssignToken, result.SemiToken},
			[]scanner.Token{
				{Type: scanner.LET, Literal: "let", LeadingTrivia: []string{" comment\n"}},
				{Type: scanner.ASSIGN, Literal: "=", LeadingTrivia: []string{"\n"}},
				{Type: scanner.SEMICOLON, Literal: ";", LeadingTrivia: []string{"c"}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
	t.Run("block", func(t *testing.T) {
		input := `
		// comment before {

		{
		let x = 100
		let y = 200 // comment before }
		/* block comment */ }`
		p := createParser(input)
		result, err := parser.ParseBlock(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]scanner.Token{result.LbraceToken, result.RbraceToken},
			[]scanner.Token{
				{Type: scanner.LBRACE, Literal: "{", LeadingTrivia: []string{"\n", " comment before {\n", "\n"}},
				{Type: scanner.RBRACE, Literal: "}", LeadingTrivia: []string{" comment before }\n", " block comment "}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
	t.Run("function", func(t *testing.T) {
		input := `/* block comment before */
	function foo
	// comment 1
	( // comment 2
	) {}`
		p := createParser(input)
		result, err := parser.ParseFunction(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]scanner.Token{result.FunctionToken, result.LparenToken, result.RparenToken},
			[]scanner.Token{
				{Type: scanner.FUNCTION, Literal: "function", LeadingTrivia: []string{" block comment before ", "\n"}},
				{Type: scanner.LPAREN, Literal: "(", LeadingTrivia: []string{"\n", " comment 1\n"}},
				{Type: scanner.RPAREN, Literal: ")", LeadingTrivia: []string{" comment 2\n"}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
	t.Run("grouped expression", func(t *testing.T) {
		input := `// comment before
	(1 + 2// comment after
	)`
		p := createParser(input)
		result, err := parser.ParseGroupedExpression(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]scanner.Token{result.LparenToken, result.RparenToken},
			[]scanner.Token{
				{Type: scanner.LPAREN, Literal: "(", LeadingTrivia: []string{" comment before\n"}},
				{Type: scanner.RPAREN, Literal: ")", LeadingTrivia: []string{" comment after\n"}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
}
