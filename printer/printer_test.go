package printer

import (
	"fmt"
	"testing"
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
						pr.PrintRune(';')
						pr.PrintRune('\n')
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
