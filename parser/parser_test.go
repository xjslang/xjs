package parser_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs/internal/debug"
	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
)

var updateGoldenFiles bool

func TestMain(m *testing.M) {
	flag.BoolVar(&updateGoldenFiles, "update", false, "update golden files")
	flag.Parse()
	os.Exit(m.Run())
}

func TestGoldenFiles(t *testing.T) {
	if updateGoldenFiles {
		t.Log("updating golden files")
	}
	files, err := filepath.Glob("./testdata/*.js")
	require.NoError(t, err)
	for _, file := range files {
		ext := filepath.Ext(file)
		goldFile := fmt.Sprintf("%s.ast.txt", strings.TrimSuffix(file, ext))
		if !updateGoldenFiles && !assert.FileExists(t, goldFile) {
			continue
		}
		// parse the source file
		source, err := os.ReadFile(file)
		require.NoError(t, err)
		s := &scanner.Scanner{}
		s.Init(source)
		p := &parser.Parser{}
		p.Init(s)
		result, err := parser.ParseProgram(p)
		require.NoError(t, err)
		// create or update golden file
		got := debug.Sprint(result)
		if updateGoldenFiles {
			err = os.WriteFile(goldFile, []byte(got), 0o644)
			require.NoError(t, err)
			continue
		}
		// compare golden file with `got`
		want, err := os.ReadFile(goldFile)
		require.NoError(t, err)
		if got != string(want) {
			diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				A:       difflib.SplitLines(got),
				B:       difflib.SplitLines(string(want)),
				Context: 5,
			})
			assert.NoError(t, err)
			t.Error(diff)
		}
	}
}

func Example_basic() {
	result, err := testutil.Parse(`function hello() {
	let x = 100
	let y = 200
}`)
	if err != nil {
		panic(err)
	}

	pr := printer.Printer{}
	pr.Init()
	pr.Print(result)
	fmt.Print(pr.String())
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}

func TestExprs(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: "1 - 2 - 3",
			expected: `BinaryExpr
	LeftValue: BinaryExpr
		LeftValue: BasicLit{Value: "1"}
		Operator: "-"
		RightValue: BasicLit{Value: "2"}
	Operator: "-"
	RightValue: BasicLit{Value: "3"}`,
		},
		{
			input: "1 + 2 * (3 + 5) - 4",
			expected: `BinaryExpr
	LeftValue: BinaryExpr
		LeftValue: BasicLit{Value: "1"}
		Operator: "+"
		RightValue: BinaryExpr
			LeftValue: BasicLit{Value: "2"}
			Operator: "*"
			RightValue: ParenExpr
				Value: BinaryExpr
					LeftValue: BasicLit{Value: "3"}
					Operator: "+"
					RightValue: BasicLit{Value: "5"}
	Operator: "-"
	RightValue: BasicLit{Value: "4"}`,
		},
		{
			input: "foo() * 2 + 1",
			expected: `BinaryExpr
	LeftValue: BinaryExpr
		LeftValue: CallExpr
			Function: Ident{Value: "foo"}
		Operator: "*"
		RightValue: BasicLit{Value: "2"}
	Operator: "+"
	RightValue: BasicLit{Value: "1"}`,
		},
		{
			input: "foo(1, 2, 3)",
			expected: `CallExpr
	Function: Ident{Value: "foo"}
	Arguments[0]: BasicLit{Value: "1"}
	Arguments[1]: BasicLit{Value: "2"}
	Arguments[2]: BasicLit{Value: "3"}`,
		},
		{
			input: "2 * (pow(2, 1 + 3) + 4)",
			expected: `BinaryExpr
	LeftValue: BasicLit{Value: "2"}
	Operator: "*"
	RightValue: ParenExpr
		Value: BinaryExpr
			LeftValue: CallExpr
				Function: Ident{Value: "pow"}
				Arguments[0]: BasicLit{Value: "2"}
				Arguments[1]: BinaryExpr
					LeftValue: BasicLit{Value: "1"}
					Operator: "+"
					RightValue: BasicLit{Value: "3"}
			Operator: "+"
			RightValue: BasicLit{Value: "4"}`,
		},
		{
			input: "1 + foo()",
			expected: `BinaryExpr
	LeftValue: BasicLit{Value: "1"}
	Operator: "+"
	RightValue: CallExpr
		Function: Ident{Value: "foo"}`,
		},
		{
			input: "1 + foo()()",
			expected: `BinaryExpr
	LeftValue: BasicLit{Value: "1"}
	Operator: "+"
	RightValue: CallExpr
		Function: CallExpr
			Function: Ident{Value: "foo"}`,
		},
	}
	for i, test := range tests {
		t.Run("exp "+strconv.Itoa(i), func(t *testing.T) {
			sc := &scanner.Scanner{}
			sc.Init([]byte(test.input))
			p := &parser.Parser{}
			p.Init(sc)
			result, err := p.ParseExpr()
			if err != nil {
				t.Fatal(err)
			}
			if got := testutil.NodeString(result); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}
