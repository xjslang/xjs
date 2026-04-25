package parser

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
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
