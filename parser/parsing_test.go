package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func createParser(input string) *parser.Parser {
	l := &scanner.Scanner{}
	l.Init([]byte(input))
	p := &parser.Parser{}
	p.Init(l)
	return p
}

func TestMalformedExpr(t *testing.T) {
	t.Run("let", func(t *testing.T) {
		input := "x = 100"
		p := createParser(input)
		_, err := parser.ParseLetStmt(p)
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
		_, err := parser.ParseFuncDecl(p)
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
			_, err := parser.ParseParenExpr(p)
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
		result, err := parser.ParseLetStmt(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.LetToken, result.AssignToken, result.SemiToken},
			[]token.Token{
				{Type: token.LET, Literal: "let", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment\n"},
				}},
				{Type: token.ASSIGN, Literal: "=", LeadingTrivia: []token.Token{
					{Type: token.NEWLINE, Literal: "\n"},
				}},
				{Type: token.SEMICOLON, Literal: ";", LeadingTrivia: []token.Token{
					{Type: token.BLOCK_COMMENT, Literal: "c"},
				}},
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
			[]token.Token{result.LbraceToken, result.RbraceToken},
			[]token.Token{
				{Type: token.LBRACE, Literal: "{", LeadingTrivia: []token.Token{
					{Type: token.NEWLINE, Literal: "\n"},
					{Type: token.LINE_COMMENT, Literal: " comment before {\n"},
					{Type: token.NEWLINE, Literal: "\n"},
				}},
				{Type: token.RBRACE, Literal: "}", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment before }\n"},
					{Type: token.BLOCK_COMMENT, Literal: " block comment "},
				}},
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
		result, err := parser.ParseFuncDecl(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.FunctionToken, result.LparenToken, result.RparenToken},
			[]token.Token{
				{Type: token.FUNCTION, Literal: "function", LeadingTrivia: []token.Token{
					{Type: token.BLOCK_COMMENT, Literal: " block comment before "},
					{Type: token.NEWLINE, Literal: "\n"},
				}},
				{Type: token.LPAREN, Literal: "(", LeadingTrivia: []token.Token{
					{Type: token.NEWLINE, Literal: "\n"},
					{Type: token.LINE_COMMENT, Literal: " comment 1\n"},
				}},
				{Type: token.RPAREN, Literal: ")", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment 2\n"},
				}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
	t.Run("grouped expression", func(t *testing.T) {
		input := `// comment before
	(1 + 2// comment after
	)`
		p := createParser(input)
		result, err := parser.ParseParenExpr(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.LparenToken, result.RparenToken},
			[]token.Token{
				{Type: token.LPAREN, Literal: "(", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment before\n"},
				}},
				{Type: token.RPAREN, Literal: ")", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment after\n"},
				}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
}

func TestParseBlockErrorRecovery(t *testing.T) {
	input := `function foo() {
		let x = 100
		let y 200 // syntax error, but keep parsing
		let z = 300
		let z = // syntax error, but keep parsing
	}

	let a = 'aaa'`
	p := createParser(input)
	_, err := parser.ParseFuncDecl(p)
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
	expected := strings.Join([]string{
		"Expected =",
		"Expected value",
	}, "\n")
	if got := err.Error(); got != expected {
		t.Fatalf("Expected %q, got %q", expected, got)
	}
	// keep parsing after `}`
	if _, err := parser.ParseLetStmt(p); err != nil {
		t.Fatal(err)
	}
}

func TestExprStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: "log()",
			expected: `ExprStmt
	Expr: CallExpr
		Function: Ident{Value: "log"}`,
		},
		{
			input: "log(1)",
			expected: `ExprStmt
	Expr: CallExpr
		Function: Ident{Value: "log"}
		Arguments[0]: BasicLit{Value: "1"}`,
		},
		{
			input: "log(1, 2)",
			expected: `ExprStmt
	Expr: CallExpr
		Function: Ident{Value: "log"}
		Arguments[0]: BasicLit{Value: "1"}
		Arguments[1]: BasicLit{Value: "2"}`,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			p := createParser(test.input)
			node, err := parser.ParseExprStmt(p)
			if err != nil {
				t.Fatal(err)
			}
			if got := testutil.NodeString(node); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}
