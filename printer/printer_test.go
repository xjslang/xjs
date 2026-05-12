package printer_test

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/printer"
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
			pr := printer.Printer{}
			pr.Init(printer.WithIndent(test.indent))
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

func TestBuiltinPrinters(t *testing.T) {
	input := `
	function foo() {
		let name = 'John Smith'
		let age = 32
		let single = true
	}
	let x = 100`
	result, err := testutil.Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	pr := printer.Printer{}
	pr.Init()
	pr.Print(result)
	fmt.Println(pr.String())
}
