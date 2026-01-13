package debug

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func TestToString(t *testing.T) {
	tests := []struct{ name, input, expected string }{
		{"LetStatement", "let x = 5", "let x=5;"},
		{"FunctionDeclaration", "function add(a, b){ return a+b }", "function add(a,b){return a+b;}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).Build(tt.input)
			stmt := p.ParseStatement()
			if stmt == nil {
				t.Fatalf("ParseStatement() returned nil")
			}
			if errors := p.Errors(); len(errors) > 0 {
				t.Fatalf("Parser errors: %v", errors)
			}
			if got := ToString(stmt); got != tt.expected {
				t.Errorf("ToString() got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	stmt := &ast.LetStatement{
		Name:  &ast.Identifier{Value: "x"},
		Value: &ast.IntegerLiteral{Token: token.Token{Literal: "5"}},
	}

	// Capture stdout to verify output
	output := captureOutput(func() {
		Print(stmt)
	})

	// Verify output is not empty
	if output == "" {
		t.Error("Print() produced no output")
	}

	// Verify output contains expected key strings
	expectedStrings := []string{"ast.LetStatement", "Identifier", "x", "IntegerLiteral", "5"}
	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Print() output missing expected string %q", expected)
		}
	}
}

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
