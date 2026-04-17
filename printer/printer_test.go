package printer

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	pr := New().WithFormat(WithSemicolon())
	pr.PrintString("hello")
	pr.PrintSemicolon()
	fmt.Println(pr.String())
}

func TestPrinter(t *testing.T) {
	tests := []struct {
		name      string
		formatted bool
		semi      bool
		indent    string
		expected  string
	}{
		{"without format", false, false, "", "block{line 0;line 1;nested block{line 0;line 1;}line 2;}"},
		{"formatted spaces", true, false, "    ", "block {\n    line 0\n    line 1\n    nested block {\n        line 0\n        line 1\n    }\n    line 2\n}"},
		{"formatted tabs", true, false, "\t", "block {\n\tline 0\n\tline 1\n\tnested block {\n\t\tline 0\n\t\tline 1\n\t}\n\tline 2\n}"},
		{"formatted tabs semi", true, true, "\t", "block {\n\tline 0;\n\tline 1;\n\tnested block {\n\t\tline 0;\n\t\tline 1;\n\t}\n\tline 2;\n}"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pr := New()
			if test.formatted {
				if test.semi {
					pr.WithFormat(WithIndent(test.indent), WithSemicolon())
				} else {
					pr.WithFormat(WithIndent(test.indent))
				}
			}
			pr.PrintString("block")
			pr.PrintWhitespace()
			pr.PrintRune('{')
			pr.PrintNewline()
			pr.IncreaseIndent()
			for i := range 3 {
				pr.PrintIndent()
				pr.PrintString(fmt.Sprintf("line %d", i))
				pr.PrintSemicolon()
				if i == 1 {
					pr.PrintNewline()
					pr.PrintIndent()
					pr.PrintString("nested block")
					pr.PrintWhitespace()
					pr.PrintRune('{')
					pr.PrintNewline()
					pr.IncreaseIndent()
					for j := range 2 {
						pr.PrintIndent()
						pr.PrintString(fmt.Sprintf("line %d", j))
						pr.PrintSemicolon()
						pr.PrintNewline()
					}
					pr.DecreaseIndent()
					pr.PrintIndent()
					pr.PrintRune('}')
				}
				pr.PrintNewline()
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
