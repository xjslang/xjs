package printer_test

import (
	"fmt"
	"strings"
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

func ExamplePrinter_Init() {
	p := &printer.Printer{}

	// Declare "middlewares" BEFORE calling Init
	p.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		if node, ok := node.(*FactorialNode); ok {
			p.Print(node.Value, "!")
			return nil
		}
		return next(node) // delegate to the "next" middleware
	})

	// Now you can use the printer
	p.Init()
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
		pr := printer.Printer{}
		pr.Init(printer.WithIndent("\t"))
		pr.Print("begin:")
		pr.IncreaseIndent()
		pr.LnPrint("aaa")
		pr.LnPrint("bbb")
		pr.LnPrint("ccc")
		pr.DecreaseIndent()
		pr.LnPrint("end")
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
			pr := printer.Printer{}
			pr.Init(printer.WithIndent(test.indent))
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
			pr := xjs.NewPrinter()
			pr.Init()
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
			pr := xjs.NewPrinter()
			pr.Init()
			js.PrintProgram(pr, node)
			out, err := pr.Output()
			require.NoError(t, err)
			if got := out; got != test.expected {
				t.Errorf("Expected\n\n%s\n\ngot\n\n%s", test.expected, got)
			}
		})
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
	out, err := pr.Output()
	require.NoError(t, err)
	if got := out; got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}

	t.Run("printing empty string does not consume ensureLine", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.Print("aaa")
		pr.EnsureLine()
		pr.Print("")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa", out)
	})

	t.Run("printing empty string does not consume ensureSpace", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.Print("aaa")
		pr.EnsureSpace()
		pr.Print("")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa", out)
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

func TestEnsureBeside(t *testing.T) {
	pr := printer.Printer{}
	pr.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		if v, ok := node.(*MyCustomStmt); ok {
			p.LnPrint(v.name)
			return nil
		}
		return next(node)
	})
	pr.Init()
	pr.Print("aaa")
	// EnsureBeside takes priority over EnsureSpace and EnsureLine
	pr.EnsureBeside()
	pr.SpPrint("bbb")
	pr.EnsureBeside()
	pr.LnPrint("bbb")
	// EnsureBeside takes priority over custom printers
	pr.BsPrint(&MyCustomStmt{name: "custom_stmt"})
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, "aaabbbbbbcustom_stmt", out)
}

func TestPrint(t *testing.T) {
	t.Run("append newlines and spaces before printing runes", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.Print("aaa")
		pr.LnPrint('b')
		pr.SpPrint('c')
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa\nb c", out)
	})
	t.Run("panic on unsupported types", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		require.Panics(t, func() { pr.Print(true) })
	})
}

func TestLnPrint(t *testing.T) {
	t.Run("new line is added before printing", func(t *testing.T) {
		pr := printer.Printer{}
		pr.Init()
		pr.LnPrint("aaa")
		pr.LnPrint("bbb")
		pr.LnPrint("ccc")
		out, err := pr.Output()
		require.NoError(t, err)
		require.Equal(t, "aaa\nbbb\nccc", out)
	})
}

func TestSpacesTakePriorityOverLines(t *testing.T) {
	pr := xjs.NewPrinter()
	pr.Print("aaa")
	// ensureSpace should take priority over ensureLine (by default statements are printed in a new line)
	pr.SpPrint(&js.ExprStmt{Expr: &js.Literal{Value: token.Token{Literal: "125"}}})
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, "aaa 125", out)
}

func TestSpPrint(t *testing.T) {
	pr := printer.Printer{}
	pr.Init()
	pr.SpPrint("aaa")
	pr.SpPrint("bbb")
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, "aaa bbb", out)
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
		{"show new lines by default", xjs.NewPrinter()},
		{"show new lines", xjs.NewPrinter(printer.WithNewLines(true))},
		{"hide new lines", xjs.NewPrinter(printer.WithNewLines(false))},
		{"hide new lines and comments", xjs.NewPrinter(printer.WithNewLines(false), printer.WithComments(false))},
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
	pr := xjs.NewPrinter(printer.Compact())
	pr.Print(result)
	out, err := pr.Output()
	require.NoError(t, err)
	golden.Assert(t, []byte(out))
}

func TestFork(t *testing.T) {
	pr := xjs.NewPrinter()
	pr.Print("for (")
	p1 := printer.Fork(pr)
	p1.LnPrint("aaa;")
	p1.SpPrint("bbb;")
	p1.SpPrint("ccc;")
	out1, err := p1.Output()
	require.NoError(t, err)
	s := out1
	s = strings.TrimRight(s, ";")
	pr.Print(s, ") {}")
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, "for (aaa; bbb; ccc) {}", out)
}
