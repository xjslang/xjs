package printer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
	"github.com/xorcare/golden"
)

type FactorialNode struct {
	Value string
}

func (node *FactorialNode) Type() string {
	return "MyCustomNode"
}

func ExamplePrinter_Init() {
	p := &printer.Printer{}

	// Declare "middlewares" BEFORE calling Init
	p.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		if node, ok := node.(*FactorialNode); ok {
			p.Print(node.Value, "!")
			return
		}
		next(node) // delegate to the "next" middleware
	})

	// Now you can use the printer
	p.Init()
	p.Print(&FactorialNode{Value: "125"})
	fmt.Println(p.String())
	// Output: 125!
}

func TestInit(t *testing.T) {
	t.Run("with custom indent", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init(printer.WithIndent("\t"))
		pr.Print("begin:")
		pr.IncreaseIndent()
		pr.LnPrint("aaa")
		pr.LnPrint("bbb")
		pr.LnPrint("ccc")
		pr.DecreaseIndent()
		pr.LnPrint("end")
		require.Equal(t, "begin:\n\taaa\n\tbbb\n\tccc\nend", pr.String())
	})
}

func TestIndent(t *testing.T) {
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
			f := printer.Printer{}
			f.Init(printer.WithIndent(test.indent))
			f.Print("block {\n")
			f.IncreaseIndent()
			for i := range 3 {
				f.PrintIndent()
				f.Print(fmt.Sprintf("line %d", i))
				f.Print(';')
				if i == 1 {
					f.Print('\n')
					f.PrintIndent()
					f.Print("nested block {\n")
					f.IncreaseIndent()
					for j := range 2 {
						f.PrintIndent()
						f.Print(fmt.Sprintf("line %d", j))
						f.Print(";\n")
					}
					f.DecreaseIndent()
					f.PrintIndent()
					f.Print('}')
				}
				f.Print('\n')
			}
			f.DecreaseIndent()
			f.PrintIndent()
			f.Print('}')
			if result := f.String(); result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestPrintCallExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "let x = 1 * (foo(2, 3 * 4) - 5) - 3",
			expected: "let x = 1 * (foo(2, 3 * 4) - 5) - 3;",
		},
		{
			input: `let x = 1 * (foo /* c1 */ (
				2 // comments before COMMA are ignored
				, // c2
				3 * 4 // c3
				// c4
			) - 5) - 3`,
			expected: `let x = 1 * (foo/* c1 */(
2, // c2
3 * 4 // c3
// c4
) - 5) - 3;`,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("exp %d", i), func(t *testing.T) {
			node, err := xjs.Parse([]byte(test.input))
			if err != nil {
				t.Fatal(err)
			}
			pr := xjs.NewPrinter()
			pr.Init()
			js.PrintProgram(pr, node)
			if got := pr.String(); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}

func TestLastComment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "in the current line",
			input:    "let x = 100 // last comment",
			expected: "let x = 100; // last comment",
		},
		{
			name:     "in a new line",
			input:    "let x = 100\n// last comment",
			expected: "let x = 100;\n// last comment",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node, err := xjs.Parse([]byte(test.input))
			if err != nil {
				t.Fatal(err)
			}
			pr := xjs.NewPrinter()
			pr.Init()
			js.PrintProgram(pr, node)
			if got := pr.String(); got != test.expected {
				t.Errorf("Expected\n\n%s\n\ngot\n\n%s", test.expected, got)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	input := "hello"
	p := printer.Printer{}
	p.Init()
	p.Print(input)
	b := p.Bytes()
	// try to modify the underlying data
	b[0] = 'H'
	expected := "hello"
	if got := p.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestEnsureLine(t *testing.T) {
	pr := printer.Printer{}
	pr.Init()
	// calling EnsureLine at the beginning of a document does not print a new line
	pr.EnsureLine()
	pr.Print("aaa")
	// calling EnsureLine multiple times only prints a new line (it is idempotent)
	for range 2 {
		pr.EnsureLine()
	}
	pr.Print("bbb")
	// calling EnsureLine on a `\n` line does not print a new line
	pr.Print("ccc\n")
	pr.EnsureLine()
	pr.Print("ddd")
	// calling EnsureLine on a `\r` line does not print a new line
	pr.Print("eee\r")
	pr.EnsureLine()
	pr.Print("fff")
	// printing empty string does not reset `ensureLine`
	pr.EnsureLine()
	pr.Print("")
	pr.Print("")
	pr.Print("ggg")
	expected := "aaa\nbbbccc\ndddeee\rfff\nggg"
	if got := pr.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}

	t.Run("printing empty string does not consume ensureLine", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.Print("aaa")
		pr.EnsureLine()
		pr.Print("")
		require.Equal(t, "aaa", pr.String())
	})

	t.Run("printing empty string does not consume ensureSpace", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.Print("aaa")
		pr.EnsureSpace()
		pr.Print("")
		require.Equal(t, "aaa", pr.String())
	})
}

func TestEnsureSpace(t *testing.T) {
	pr := printer.Printer{}
	pr.Init()
	// calling EnsureSpace at the beginning of a document does not print a new space
	pr.EnsureSpace()
	pr.Print("aaa")
	// calling EnsureSpace multiple times only prints a new space (it is idempotent)
	for range 2 {
		pr.EnsureSpace()
	}
	pr.Print("bbb")
	// calling EnsureSpace on a `\n` line does not print a new space
	pr.Print("ccc\n")
	pr.EnsureSpace()
	pr.Print("ddd")
	// calling EnsureSpace on a `\r` line does not print a new space
	pr.Print("eee\r")
	pr.EnsureSpace()
	pr.Print("fff")
	// printing empty string does not reset `ensureSpace`
	pr.EnsureSpace()
	pr.Print("")
	pr.Print("")
	pr.Print("ggg")
	expected := "aaa bbbccc\ndddeee\rfff ggg"
	if got := pr.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestPrint(t *testing.T) {
	t.Run("append newlines and spaces before printing runes", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.Print("aaa")
		pr.LnPrint('b')
		pr.SpPrint('c')
		require.Equal(t, "aaa\nb c", pr.String())
	})
	t.Run("panic on unsupported types", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		require.Panics(t, func() { pr.Print(100) })
	})
}

func TestLnPrint(t *testing.T) {
	t.Run("new line is added before printing", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.LnPrint("aaa")
		pr.LnPrint("bbb")
		pr.LnPrint("ccc")
		require.Equal(t, "aaa\nbbb\nccc", pr.String())
	})
}

func TestSpPrint(t *testing.T) {
	t.Run("spaces are added in between", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.SpPrint("aaa", "bbb", "ccc")
		require.Equal(t, "aaa bbb ccc", pr.String())
	})
}

func TestWithComments(t *testing.T) {
	input := "// c\nlet x = 100\n/* c */let y = 200"
	tests := []struct {
		name string
		pr   *printer.Printer
	}{
		{"show comments by default", xjs.NewPrinter()},
		{"hide comments", xjs.NewPrinter(printer.WithComments(false))},
		{"show comments", xjs.NewPrinter(printer.WithComments(true))},
		{"hide lcomments and show bcomments", xjs.NewPrinter(printer.WithLineComments(false), printer.WithBlockComments(true))},
		{"show lcomments and hide bcomments", xjs.NewPrinter(printer.WithLineComments(true), printer.WithBlockComments(false))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := xjs.Parse([]byte(input))
			require.NoError(t, err)
			test.pr.Print(result)
			golden.Assert(t, test.pr.Bytes())
		})
	}
}

func TestWithEmptyLines(t *testing.T) {
	input := "let x = 100\n\n\n// line comment\nlet y = 200"
	tests := []struct {
		name string
		pr   *printer.Printer
	}{
		{"show empty lines by default", xjs.NewPrinter()},
		{"show empty lines", xjs.NewPrinter(printer.WithEmptyLines(true))},
		{"hide empty lines", xjs.NewPrinter(printer.WithEmptyLines(false))},
		{"hide empty lines and comments", xjs.NewPrinter(printer.WithEmptyLines(false), printer.WithComments(false))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := xjs.Parse([]byte(input))
			require.NoError(t, err)
			test.pr.Print(result)
			golden.Assert(t, test.pr.Bytes())
		})
	}
}
