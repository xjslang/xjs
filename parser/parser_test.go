package parser

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/source"
)

func ExampleParse() {
	result, err := Parse([]byte(`function hello() {
	let x = 100
	let y = 200
}`))
	if err != nil {
		panic(err)
	}

	pr := printer.New()
	result.PrintTo(pr)
	fmt.Print(pr.String())
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single value", "3", `*ast.IntegerLiteral{Value: "3"}`},
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

func TestParser(t *testing.T) {
	input := `
		function printHello() {
			let x = 100;
			let y = 200;
		}

		let x = 100;
		let y = 200;`
	l := &lexer.Lexer{}
	l.Init([]byte(input))
	p := Parser{}
	p.Init(l)
	pr, err := p.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}
	expectNode(t, pr, &ast.BlockStatement{
		Statements: []ast.Statement{
			&ast.FunctionDeclaration{
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.LetStatement{},
						&ast.LetStatement{},
					},
				},
			},
			&ast.LetStatement{},
			&ast.LetStatement{},
		},
	})
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Error
	}{
		{"missing token", "function\nhello({}", Error{
			Range: source.Range{
				Start: source.Position{Line: 1, Column: 6},
				End:   source.Position{Line: 1, Column: 7},
			},
			Message: "Expected )"},
		},
		{"missing semicolon", "let x = 100 let y = 200", Error{
			Range: source.Range{
				Start: source.Position{Line: 0, Column: 12},
				End:   source.Position{Line: 0, Column: 15},
			},
			Message: "Expected statement terminator"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer.Lexer{}
			l.Init([]byte(test.input))
			p := Parser{}
			p.Init(l)
			_, err := p.ParseProgram()
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			list, ok := err.(ErrorList)
			if !ok {
				t.Fatalf("Expected ErrorList, got %T", err)
			}
			if n := len(list); n != 1 {
				t.Errorf("Expected one error, got %d", n)
			} else if list[0].Message != test.expected.Message {
				t.Errorf("Expected error %q, got %q", test.expected, list[0].Message)
			} else if list[0].Range != test.expected.Range {
				t.Error("position error")
			}
		})
	}

	t.Run("multiple errors", func(t *testing.T) {
		expectedErrors := []string{
			"Expected value",
			"Expected =",
			"Expected statement terminator",
			"Expected )",
		}

		l := &lexer.Lexer{}
		l.Init([]byte(`
		let x = !
		let y
		let x = 100.

		function print(`))
		p := Parser{}
		p.Init(l)
		_, err := p.ParseProgram()
		errors, ok := err.(ErrorList)
		if !ok {
			t.Fatalf("Expected ErrorList, got %T", err)
		}
		if n := len(errors); n != len(expectedErrors) {
			t.Fatalf("Expected %d errors, got %d", len(expectedErrors), n)
		}
		for i, expectedError := range expectedErrors {
			if errors[i].Message != expectedError {
				t.Errorf("Expected %q, got %q", expectedError, errors[i].Message)
			}
		}
	})
}

func expectNode(t *testing.T, a ast.Statement, b ast.Statement) {
	switch expected := b.(type) {
	case *ast.BlockStatement:
		got, ok := a.(*ast.BlockStatement)
		if !ok {
			t.Errorf("Expected %T, got %T", b, a)
		}
		if l := len(expected.Statements); l != len(got.Statements) {
			t.Errorf("Expected %d statements, got %d", l, len(got.Statements))
			return
		}
		for i, expectedStmt := range expected.Statements {
			expectNode(t, got.Statements[i], expectedStmt)
		}
	case *ast.FunctionDeclaration:
		got, ok := a.(*ast.FunctionDeclaration)
		if !ok {
			t.Errorf("Expected %T, got %T", b, a)
		}
		expectNode(t, got.Body, expected.Body)
	case *ast.LetStatement:
		_, ok := a.(*ast.LetStatement)
		if !ok {
			t.Errorf("Expected %T, got %T", b, a)
		}
	}
}
