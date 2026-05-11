package printer

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/js"
)

func TestPrinter(t *testing.T) {
	tests := []struct {
		name     string
		indent   string
		expected string
	}{
		{"with spaces", "    ", "block {\n    line 0;\n    line 1;\n    nested block {\n        line 0;\n        line 1;\n    }\n    line 2;\n}"},
		{"with tabs", "\t", "block {\n\tline 0;\n\tline 1;\n\tnested block {\n\t\tline 0;\n\t\tline 1;\n\t}\n\tline 2;\n}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pr := New(WithIndent(test.indent))
			pr.PrintString("block {\n")
			pr.IncreaseIndent()
			for i := range 3 {
				pr.PrintIndent()
				pr.PrintString(fmt.Sprintf("line %d", i))
				pr.PrintRune(';')
				if i == 1 {
					pr.PrintRune('\n')
					pr.PrintIndent()
					pr.PrintString("nested block {\n")
					pr.IncreaseIndent()
					for j := range 2 {
						pr.PrintIndent()
						pr.PrintString(fmt.Sprintf("line %d", j))
						pr.PrintString(";\n")
					}
					pr.DecreaseIndent()
					pr.PrintIndent()
					pr.PrintRune('}')
				}
				pr.PrintRune('\n')
			}
			pr.DecreaseIndent()
			pr.PrintIndent()
			pr.PrintRune('}')
			if result := pr.String(); result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestIncreaseDecreaseIndent(t *testing.T) {
	t.Run("indent level changes correctly", func(t *testing.T) {
		pr := New()
		// at first indentLevel must be 0
		if pr.indentLevel != 0 {
			t.Errorf("Expected indentLevel to be 0, got %d", pr.indentLevel)
		}
		// increase indentLevel
		pr.IncreaseIndent()
		if pr.indentLevel != 1 {
			t.Errorf("Expected indentLevel to be 1, got %d", pr.indentLevel)
		}
		// increase indentLevel twice
		pr.IncreaseIndent()
		pr.IncreaseIndent()
		if pr.indentLevel != 3 {
			t.Errorf("Expected indentLevel to be 3, got %d", pr.indentLevel)
		}
		// decrease indentLevel
		pr.DecreaseIndent()
		if pr.indentLevel != 2 {
			t.Errorf("Expected indentLevel to be 2, got %d", pr.indentLevel)
		}
		// decrease indentLevel twice
		pr.DecreaseIndent()
		pr.DecreaseIndent()
		if pr.indentLevel != 0 {
			t.Errorf("Expected indentLevel to be 0, got %d", pr.indentLevel)
		}
	})
	t.Run("cannot decrease below zero", func(t *testing.T) {
		pr := New()
		// indentLevel cannot be negative
		pr.DecreaseIndent()
		if pr.indentLevel != 0 {
			t.Errorf("Expected indentLevel to be 0, got %d", pr.indentLevel)
		}
	})
}

func TestUsePrinter(t *testing.T) {
	input := `
	function foo() {
		let a = x * (200 + 3)
		let b = 200
	}
	let x = 100
	let y = 200
	let z = 'aaa'
	let v = true
	let w = false
`
	result, err := js.Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	c := Printer{}
	c.UsePrinter(IntegerLiteralPrinter)
	c.UsePrinter(StringLiteralPrinter)
	c.UsePrinter(BooleanLiteralPrinter)
	c.UsePrinter(InfixOperatorPrinter)
	c.UsePrinter(GroupedExpressionPrinter)
	c.UsePrinter(IdentifierPrinter)
	c.UsePrinter(LetPrinter)
	c.UsePrinter(FunctionPrinter)
	c.UsePrinter(BlockPrinter)
	c.Print(result)
	fmt.Println(c.String())
}
