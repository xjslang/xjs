package parser

import (
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/lexer"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single value", "3", `*ast.IntegerLiteral{Value: "3"}`},
		{
			name:  "boolean values",
			input: "true + false",
			expected: `*ast.InfixOperator
	LeftValue: *ast.BooleanLiteral{Value: "true"}
	Operator: "+"
	RightValue: *ast.BooleanLiteral{Value: "false"}`,
		},
		{
			name:  "identifiers",
			input: "x + 2*y",
			expected: `*ast.InfixOperator
	LeftValue: *ast.Identifier{Value: "x"}
	Operator: "+"
	RightValue: *ast.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "2"}
		Operator: "*"
		RightValue: *ast.Identifier{Value: "y"}`,
		},
		{
			name:  "basic",
			input: "1 + 2 - 3",
			expected: `*ast.InfixOperator
	LeftValue: *ast.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "1"}
		Operator: "+"
		RightValue: *ast.IntegerLiteral{Value: "2"}
	Operator: "-"
	RightValue: *ast.IntegerLiteral{Value: "3"}`,
		},
		{
			name:  "complex",
			input: "5 - 2 * '3' % 5 + 1",
			expected: `*ast.InfixOperator
	LeftValue: *ast.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "5"}
		Operator: "-"
		RightValue: *ast.InfixOperator
			LeftValue: *ast.InfixOperator
				LeftValue: *ast.IntegerLiteral{Value: "2"}
				Operator: "*"
				RightValue: *ast.StringLiteral{Value: "'3'"}
			Operator: "%"
			RightValue: *ast.IntegerLiteral{Value: "5"}
	Operator: "+"
	RightValue: *ast.IntegerLiteral{Value: "1"}`,
		},
		{
			name:  "parentheses",
			input: "2 * (3 + 5) - (1)",
			expected: `*ast.InfixOperator
	LeftValue: *ast.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "2"}
		Operator: "*"
		RightValue: *ast.GroupedExpression
			Value: *ast.InfixOperator
				LeftValue: *ast.IntegerLiteral{Value: "3"}
				Operator: "+"
				RightValue: *ast.IntegerLiteral{Value: "5"}
	Operator: "-"
	RightValue: *ast.GroupedExpression
		Value: *ast.IntegerLiteral{Value: "1"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer.Lexer{}
			l.Init([]byte(test.input))
			p := Parser{}
			p.Init(l)
			exp, err := p.ParseExpression()
			if err != nil {
				t.Fatal(err)
			}
			if result := testutil.NodeString(exp); result != test.expected {
				t.Errorf("Invalid node:\nExpected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

func TestMalformedExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "missing value",
			input:    "1 + 2*",
			expected: "Expected value",
		},
		{
			name:     "missing right parenthesis",
			input:    "2 * (3 + 5",
			expected: "Expected )",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer.Lexer{}
			l.Init([]byte(test.input))
			p := Parser{}
			p.Init(l)
			_, err := p.ParseExpression()
			if err == nil {
				t.Fatal("An error was expected")
			}
			if result := err.Error(); result != test.expected {
				t.Errorf("Expected error to be %q, got %q", test.expected, result)
			}
		})
	}
}
