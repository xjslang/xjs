package parser

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/printer"
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

func TestParser(t *testing.T) {
	input := `
		function printHello() {
			let x = 100;
			let y = 200;
		}

		let x = 100;
		let y = 200;`
	l := lexer.New([]byte(input))
	p := newParser(l)
	pr, err := p.parseProgram()
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
		expected string
	}{
		{"missing token", "function hello({}", "Expected RPAREN, got LBRACE"},
		{"missing semicolon", "let x = 100 let y = 200", "Expected semicolon, newline, or EOF, got LET"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.New([]byte(test.input))
			p := newParser(l)
			_, err := p.parseProgram()
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			list, ok := err.(ErrorList)
			if !ok {
				t.Fatalf("Expected ErrorList, got %T", err)
			}
			if n := len(list); n != 1 {
				t.Errorf("Expected one error, got %d", n)
			} else if list[0] != test.expected {
				t.Errorf("Expected error %q, got %q", test.expected, list[0])
			}
		})
	}
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
