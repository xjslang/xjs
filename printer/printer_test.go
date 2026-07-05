package printer_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
	"github.com/xorcare/golden"
)

type FactorialNode struct {
	ast.BaseExpr
	Value string
}

func ExampleBuilder_Build() {
	p := printer.NewBuilder().
		UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
			if node, ok := node.(*FactorialNode); ok {
				p.Print(node.Value, "!")
				return nil
			}
			return next(node) // delegate to the "next" middleware
		}).
		Build()
	p.Print(&FactorialNode{Value: "125"})
	out, err := p.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: 125!
}

func TestInit(t *testing.T) {
	t.Run("with custom indent", func(t *testing.T) {
		pr := printer.NewBuilder().Build(printer.WithIndent("\t"))
		pr.Print("begin:")
		pr.IncreaseIndent()
		pr.Line().Print("aaa")
		pr.Line().Print("bbb")
		pr.Line().Print("ccc")
		pr.DecreaseIndent()
		pr.Line().Print("end")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "begin:\n\taaa\n\tbbb\n\tccc\nend", out)
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
			pr := printer.NewBuilder().Build(printer.WithIndent(test.indent))
			pr.Print("block {\n")
			pr.IncreaseIndent()
			for i := range 3 {
				pr.PrintIndent()
				pr.Print(fmt.Sprintf("line %d", i))
				pr.Print(';')
				if i == 1 {
					pr.Print('\n')
					pr.PrintIndent()
					pr.Print("nested block {\n")
					pr.IncreaseIndent()
					for j := range 2 {
						pr.PrintIndent()
						pr.Print(fmt.Sprintf("line %d", j))
						pr.Print(";\n")
					}
					pr.DecreaseIndent()
					pr.PrintIndent()
					pr.Print('}')
				}
				pr.Print('\n')
			}
			pr.DecreaseIndent()
			pr.PrintIndent()
			pr.Print('}')
			out, err := pr.Output()
			require.NoError(t, err)
			if result := out; result != test.expected {
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
			expected: `let x = 1 * (foo /* c1 */(
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
			pr := xjs.PrinterBuilder().Build()
			js.PrintProgram(pr, node)
			out, err := pr.Output()
			require.NoError(t, err)
			if got := out; got != test.expected {
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
			pr := xjs.PrinterBuilder().Build()
			js.PrintProgram(pr, node)
			out, err := pr.Output()
			require.NoError(t, err)
			if got := out; got != test.expected {
				t.Errorf("Expected\n\n%s\n\ngot\n\n%s", test.expected, got)
			}
		})
	}
}

func TestLine(t *testing.T) {
	pr := printer.NewBuilder().Build()
	// calling Line at the beginning of a document does not print a new line
	pr.Line()
	pr.Print("aaa")
	// calling Line multiple times only prints a new line (it is idempotent)
	for range 2 {
		pr.Line()
	}
	pr.Print("bbb")
	// calling Line on a `\n` line does not print a new line
	pr.Print("ccc\n")
	pr.Line()
	pr.Print("ddd")
	// calling Line on a `\r` line does not print a new line
	pr.Print("eee\r")
	pr.Line()
	pr.Print("fff")
	// printing empty string does not reset `ensureLine`
	pr.Line()
	pr.Print("")
	pr.Print("")
	pr.Print("ggg")
	expected := "aaa\nbbbccc\ndddeee\rfff\nggg"
	out, err := pr.Output()
	require.NoError(t, err)
	if got := out; got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}

	t.Run("printing empty string does not consume ensureLine", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		pr.Print("aaa")
		pr.Line()
		pr.Print("")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa", out)
	})

	t.Run("printing empty string does not consume ensureSpace", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		pr.Print("aaa")
		pr.Space()
		pr.Print("")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa", out)
	})
}

func TestSpace(t *testing.T) {
	pr := printer.NewBuilder().Build()
	// calling Space at the beginning of a document does not print a new space
	pr.Space()
	pr.Print("aaa")
	// calling Space multiple times only prints a new space (it is idempotent)
	for range 2 {
		pr.Space()
	}
	pr.Print("bbb")
	// calling Space on a `\n` line does not print a new space
	pr.Print("ccc\n")
	pr.Space()
	pr.Print("ddd")
	// calling Space on a `\r` line does not print a new space
	pr.Print("eee\r")
	pr.Space()
	pr.Print("fff")
	// printing empty string does not reset `ensureSpace`
	pr.Space()
	pr.Print("")
	pr.Print("")
	pr.Print("ggg")
	expected := "aaa bbbccc\ndddeee\rfff ggg"
	out, err := pr.Output()
	require.NoError(t, err)
	if got := out; got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

type MyCustomStmt struct {
	ast.BaseStmt
	name string
}

func TestBeside(t *testing.T) {
	pr := printer.NewBuilder().
		UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
			if v, ok := node.(*MyCustomStmt); ok {
				p.Line().Print(v.name)
				return nil
			}
			return next(node)
		}).
		Build()
	pr.Print("aaa")
	// Beside takes priority over Space and Line
	pr.Beside()
	pr.Space().Print("bbb")
	pr.Beside()
	pr.Line().Print("bbb")
	// Beside takes priority over custom printers
	pr.Beside().Print(&MyCustomStmt{name: "custom_stmt"})
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, "aaabbbbbbcustom_stmt", out)
}

func TestPrint(t *testing.T) {
	t.Run("append newlines and spaces before printing runes", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		pr.Print("aaa")
		pr.Line().Print('b')
		pr.Space().Print('c')
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa\nb c", out)
	})
	t.Run("panic on unsupported types", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		require.Panics(t, func() { pr.Print(true) })
	})
}

func TestLineAndPrint(t *testing.T) {
	t.Run("new line is added before printing", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		pr.Line().Print("aaa")
		pr.Line().Print("bbb")
		pr.Line().Print("ccc")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa\nbbb\nccc", out)
	})
}

func TestSpacesTakePriorityOverLines(t *testing.T) {
	pr := xjs.PrinterBuilder().Build()
	pr.Print("aaa")
	// ensureSpace should take priority over ensureLine (by default statements are printed in a new line)
	pr.Space().Print(&js.ExprStmt{Expr: &js.Literal{Value: token.Token{Literal: "125"}}})
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, "aaa 125", out)
}

func TestSpaceAndPrint(t *testing.T) {
	pr := printer.NewBuilder().Build()
	pr.Space().Print("aaa")
	pr.Space().Print("bbb")
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, "aaa bbb", out)
}

func TestPrintPriority(t *testing.T) {
	t.Run("Beside takes priority over Space and Line", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		pr.Print("a").
			Beside().Space().Print("b").
			Beside().Line().Print("c")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "abc", out)
	})
	t.Run("Space takes priority over Line", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		pr.Print("a").
			Space().Line().Print("b")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "a b", out)
	})
	t.Run("Line has the lowest priority", func(t *testing.T) {
		pr := printer.NewBuilder().Build()
		pr.Print("a").
			Line().Print("b")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "a\nb", out)
	})
}

func TestWithComments(t *testing.T) {
	input := "// c\nlet x = 100\n/* c */let y = 200"
	tests := []struct {
		name string
		pr   *printer.Printer
	}{
		{"show comments by default", xjs.PrinterBuilder().Build()},
		{"hide comments", xjs.PrinterBuilder().Build(printer.WithComments(false))},
		{"show comments", xjs.PrinterBuilder().Build(printer.WithComments(true))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := xjs.Parse([]byte(input))
			require.NoError(t, err)
			test.pr.Print(result)
			out, err := test.pr.Output()
			require.NoError(t, err)
			golden.Assert(t, []byte(out))
		})
	}
}

func TestWithNewLines(t *testing.T) {
	input := "let x = 100\n\n\n// line comment\nlet y = 200"
	tests := []struct {
		name string
		pr   *printer.Printer
	}{
		{"show new lines by default", xjs.PrinterBuilder().Build()},
		{"show new lines", xjs.PrinterBuilder().Build(printer.WithNewLines(true))},
		{"hide new lines", xjs.PrinterBuilder().Build(printer.WithNewLines(false))},
		{"hide new lines and comments", xjs.PrinterBuilder().Build(printer.WithNewLines(false), printer.WithComments(false))},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := xjs.Parse([]byte(input))
			require.NoError(t, err)
			test.pr.Print(result)
			out, err := test.pr.Output()
			require.NoError(t, err)
			golden.Assert(t, []byte(out))
		})
	}
}

func TestCompact(t *testing.T) {
	input := `
	// c
	function foo() {
		let a = 'aaa'; // c
		let b = 'bbb'; // c
	}
	// c
	let x = 100
	/* b */
	let y = 200`
	result, err := xjs.Parse([]byte(input))
	require.NoError(t, err)
	out, err := xjs.Print(result, printer.Compact())
	require.NoError(t, err)
	golden.Assert(t, []byte(out))
}
