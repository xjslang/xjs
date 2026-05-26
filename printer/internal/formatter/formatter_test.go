package formatter_test

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/printer/internal/formatter"
)

func TestString(t *testing.T) {
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
			f := formatter.Formatter{}
			f.Init(formatter.WithIndent(test.indent))
			f.PrintString("block {\n")
			f.IncreaseIndent()
			for i := range 3 {
				f.PrintIndent()
				f.PrintString(fmt.Sprintf("line %d", i))
				f.PrintRune(';')
				if i == 1 {
					f.PrintRune('\n')
					f.PrintIndent()
					f.PrintString("nested block {\n")
					f.IncreaseIndent()
					for j := range 2 {
						f.PrintIndent()
						f.PrintString(fmt.Sprintf("line %d", j))
						f.PrintString(";\n")
					}
					f.DecreaseIndent()
					f.PrintIndent()
					f.PrintRune('}')
				}
				f.PrintRune('\n')
			}
			f.DecreaseIndent()
			f.PrintIndent()
			f.PrintRune('}')
			if result := f.String(); result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}
