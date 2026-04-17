package ast

import (
	"testing"

	"github.com/xjslang/xjs/printer"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{"identifier", &Identifier{Value: "hello"}, "hello"},
		{"integer literal", &IntegerLiteral{Value: "12345"}, "12345"},
		{"string literal", &StringLiteral{Value: "'Hello, World!'"}, "'Hello, World!'"},
		{"let statement", &LetStatement{
			Name:  &Identifier{Value: "x"},
			Value: &IntegerLiteral{Value: "125"},
		}, "let x = 125"},
		{"function declaration", &FunctionDeclaration{
			Name: &Identifier{Value: "foo"},
			Body: &BlockStatement{
				Statements: []Statement{
					&LetStatement{Name: &Identifier{Value: "x"}, Value: &IntegerLiteral{Value: "100"}},
					&LetStatement{Name: &Identifier{Value: "y"}, Value: &IntegerLiteral{Value: "200"}},
				},
			},
		}, "function foo() {let x = 100;let y = 200}"},
	}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pr := &printer.Printer{}
			test.node.PrintTo(pr)
			if s := pr.String(); s != test.expected {
				t.Errorf("node %d: expected %q, got %q", i, test.expected, s)
			}
		})
	}
}
