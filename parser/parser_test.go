package parser

import (
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func TestParser(t *testing.T) {
	input := `
		function printHello() {
			let x = 100;
			let y = 200;
		}

		let x = 100;
		let y = 200;`
	l := lexer.New([]byte(input))
	p := New(l)
	pr := p.ParseProgram()
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

// Test Automatic Semicolon Insertion (ASI for short)
func TestASI(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"newline in the middle", "let x = 100\nlet y = 200;"},
		{"newline at the end", "let x = 100\n"},
		{"eof at the end", "let x = 100"},
		{"newline in block comment", "let x = 100/* block\ncomment */let y = 200;"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := lexer.New([]byte(test.input))
			p := New(l)
			if _, err := p.parseLetStatement(); err != nil {
				t.Error(err)
			}
		})
	}

	t.Run("block comment without newlines", func(t *testing.T) {
		l := lexer.New([]byte("let x = 100/* block comment */let y = 200"))
		p := New(l)
		if _, err := p.parseLetStatement(); err == nil {
			t.Error("Expected ASI error, nothing happened")
		}
	})

	t.Run("current token not advanced on implicit semicolon", func(t *testing.T) {
		l := lexer.New([]byte("let x = 100\nlet y = 200"))
		p := New(l)
		_, err := p.parseLetStatement()
		if err != nil {
			t.Fatal(err)
		}
		if p.CurrentToken.Type != token.LET {
			t.Errorf("Expected %v, got %v", token.LET, p.CurrentToken.Type)
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
