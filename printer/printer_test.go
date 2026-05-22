package printer_test

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/printer"
)

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
			node, err := testutil.Parse(test.input)
			if err != nil {
				t.Fatal(err)
			}
			pr := &printer.Printer{}
			pr.Init()
			printer.PrintProgram(pr, node)
			if got := pr.String(); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}

func TestCommentFormatting(t *testing.T) {
	input := `//a
		/*b*/ function//c
						foo/*c*/(//d
						)//e
	{
						//f
						function boo(){
							let a = 'aaa';/*g*/let b = 'bbb'/*h*/;
						}
		//c
	}
	let x = 100 //y
	let y = (/*c*/100 //j
	+//k
	200)
	log(x, y)/*l*/;
	let z = 1
	+ 2 // last comment`
	expected := `//a
/*b*/
function //c
foo/*c*/( //d
) //e
{
  //f
  function boo() {
    let a = 'aaa';/*g*/
    let b = 'bbb'/*h*/;
  }
  //c
}
let x = 100; //y
let y = (/*c*/100 //j
+ //k
200);
log(x, y)/*l*/;
let z = 1
+ 2; // last comment`
	node, err := testutil.Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	pr := &printer.Printer{}
	pr.Init()
	printer.PrintProgram(pr, node)
	if got := pr.String(); got != expected {
		t.Errorf("Expected\n\n%s\n\ngot\n\n%s", expected, got)
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
			node, err := testutil.Parse(test.input)
			if err != nil {
				t.Fatal(err)
			}
			pr := &printer.Printer{}
			pr.Init()
			printer.PrintProgram(pr, node)
			if got := pr.String(); got != test.expected {
				t.Errorf("Expected\n\n%s\n\ngot\n\n%s", test.expected, got)
			}
		})
	}
}

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

func TestBytes(t *testing.T) {
	input := "hello"
	p := printer.Printer{}
	p.Init()
	p.PrintString(input)
	b := p.Bytes()
	// try to modify the underlaying data
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
	pr.PrintIndentedString("aaa")
	// calling EnsureLine multiple times only prints a new line (it is idempotent)
	for range 2 {
		pr.EnsureLine()
	}
	pr.PrintIndentedString("bbb")
	// calling EnsureLine on a `\n` line does not print a new line
	pr.PrintIndentedString("ccc\n")
	pr.EnsureLine()
	pr.PrintIndentedString("ddd")
	// calling EnsureLine on a `\r` line does not print a new line
	pr.PrintIndentedString("eee\r")
	pr.EnsureLine()
	pr.PrintIndentedString("fff")
	// printing empty string does not reset `ensureLine`
	pr.EnsureLine()
	pr.PrintString("")
	pr.PrintIndentedString("")
	pr.PrintIndentedString("ggg")
	expected := "aaa\nbbbccc\ndddeee\rfff\nggg"
	if got := pr.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func TestEnsureSpace(t *testing.T) {
	pr := printer.Printer{}
	pr.Init()
	// calling EnsureSpace at the beginning of a document does not print a new space
	pr.EnsureSpace()
	pr.PrintString("aaa")
	// calling EnsureSpace multiple times only prints a new space (it is idempotent)
	for range 2 {
		pr.EnsureSpace()
	}
	pr.PrintIndentedString("bbb")
	// calling EnsureSpace on a `\n` line does not print a new space
	pr.PrintIndentedString("ccc\n")
	pr.EnsureSpace()
	pr.PrintIndentedString("ddd")
	// calling EnsureSpace on a `\r` line does not print a new space
	pr.PrintIndentedString("eee\r")
	pr.EnsureSpace()
	pr.PrintIndentedString("fff")
	// printing empty string does not reset `ensureSpace`
	pr.EnsureSpace()
	pr.PrintString("")
	pr.PrintIndentedString("")
	pr.PrintIndentedString("ggg")
	expected := "aaa bbbccc\ndddeee\rfff ggg"
	if got := pr.String(); got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
