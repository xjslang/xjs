package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/token"
)

func TestMalformedExpr(t *testing.T) {
	t.Run("block", func(t *testing.T) {
		tests := []struct {
			input       string
			expectedErr string
		}{
			{"let x = 100 }", "Expected {"},
			{"{ let x = 100", "Expected }"},
		}
		for i, test := range tests {
			p := xjs.NewBuilder().Build([]byte(test.input))
			_, err := js.ParseBlockStmt(p)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			if got := err.Error(); got != test.expectedErr {
				t.Fatalf("%d: Expected error to be %q, got %q", i, test.expectedErr, got)
			}
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
			p := xjs.NewBuilder().Build([]byte(test.input))
			_, err := js.ParseParenExpr(p)
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
	t.Run("block", func(t *testing.T) {
		input := `
		// comment before {

		{
		let x = 100
		let y = 200 // comment before }
		/* block comment */ }`
		p := xjs.NewBuilder().Build([]byte(input))
		result, err := js.ParseBlockStmt(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.Tokens.Lbrace, result.Tokens.Rbrace},
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
	t.Run("grouped expression", func(t *testing.T) {
		input := `// comment before
	(1 + 2// comment after
	)`
		p := xjs.NewBuilder().Build([]byte(input))
		result, err := js.ParseParenExpr(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.Tokens.Lparen, result.Tokens.Rparen},
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

func TestStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: "log()",
			expected: `*js.ExprStmt
	Expr: *js.CallExpr
		Callee: *js.Variable{Name: "log"}`,
		},
		{
			input: "log(1)",
			expected: `*js.ExprStmt
	Expr: *js.CallExpr
		Callee: *js.Variable{Name: "log"}
		Args[0]: *js.Literal{Value: "1"}`,
		},
		{
			input: "log(1, 2)",
			expected: `*js.ExprStmt
	Expr: *js.CallExpr
		Callee: *js.Variable{Name: "log"}
		Args[0]: *js.Literal{Value: "1"}
		Args[1]: *js.Literal{Value: "2"}`,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			p := xjs.NewBuilder().Build([]byte(test.input))
			node, err := js.ParseExprStmt(p)
			if err != nil {
				t.Fatal(err)
			}
			if got := testutil.NodeString(node); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}

func TestInvalidTokenAfterNewline(t *testing.T) {
	tests := []struct {
		input  string
		errors []string
	}{
		{"\n%", []string{"Expected value"}},
		{"let\n%", []string{"Expected identifier", "Expected value"}},
		{"let x\n%", []string{"Expected =", "Expected value"}},
		{"let y =\n%", []string{"Expected value", "Expected value"}},
		{"let x =\nlet y = 1", []string{"Expected value"}},
	}
	for i := range 2 {
		for j, test := range tests {
			t.Run(fmt.Sprintf("test %d%d", i, j), func(t *testing.T) {
				var input string
				if i > 0 {
					input = fmt.Sprintf("{%s}", test.input)
				} else {
					input = test.input
				}
				p := xjs.NewBuilder().Build([]byte(input))
				var err error
				if i > 0 {
					_, err = js.ParseBlockStmt(p)
				} else {
					_, err = js.ParseProgram(p)
				}
				if err == nil {
					t.Fatal("Expected an error, got nil")
				}
				expected := strings.Join(test.errors, "\n")
				if got := err.Error(); got != expected {
					t.Errorf("Expected %q, got %q", expected, got)
				}
			})
		}
	}
}
