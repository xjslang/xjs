package xjs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

func Example_basic() {
	input := `function hello() {
	let x = 100
	let y = 200
}`
	result, err := xjs.Parse([]byte(input))
	if err != nil {
		panic(err)
	}

	pr := xjs.NewPrinter()
	pr.Init()
	pr.Print(result)
	fmt.Print(pr.String())
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}

func TestLanguageFeatures(t *testing.T) {
	pattern := filepath.Join("testdata", "*.js")
	files, err := filepath.Glob(pattern)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		testName := strings.TrimSuffix(filepath.Base(file), ".js")
		t.Run(testName, func(t *testing.T) {
			// read and parse file
			dat, err := os.ReadFile(file)
			require.NoError(t, err)
			result, err := xjs.Parse(dat)
			require.NoError(t, err)
			// print without newlines trivia and parse it again
			pr := xjs.NewPrinter()
			pr.Print(result)
			out, _ := pr.Output()
			result, err = xjs.Parse([]byte(out))
			require.NoError(t, err)
			// print with default options
			pr = xjs.NewPrinter()
			pr.Print(result)
			// the original must match the final printed result
			require.Equal(t, string(dat), pr.String())
		})
	}
}

func TestParseCommaDangle(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{
			input:    "let a = {\n\tone: 'one',\n\ttwo: 'two',\n}",
			expected: "let a = {\n\tone: 'one',\n\ttwo: 'two'\n};",
		},
		{
			input:    "let a = [\n\t'one',\n\t'two',\n]",
			expected: "let a = [\n\t'one',\n\t'two'\n];",
		},
		{
			input:    "let a = point(\n\tx, \n\ty,\n)",
			expected: "let a = point(\n\tx,\n\ty\n);",
		},
		{
			input:    "let a = function (\n\tx, \n\ty,\n) {}",
			expected: "let a = function (\n\tx,\n\ty\n) {};",
		},
		{
			input:    "function point(\n\tx, \n\ty,\n) {}",
			expected: "function point(\n\tx,\n\ty\n) {}",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			result, err := xjs.Parse([]byte(test.input))
			require.NoError(t, err)
			pr := xjs.NewPrinter(printer.WithIndent("\t"))
			pr.Print(result)
			assert.Equal(t, test.expected, pr.String())
		})
	}
}

type iifeExpr struct {
	ast.BaseExpr
	LparenToken token.Token
	RparenToken token.Token
	Function    *js.FunctionDecl
}

func TestMiddlewares(t *testing.T) {
	input := `(function foo() {
		print('Hello, World!')
	})()`
	b := xjs.NewBuilder()
	// parse IIFE expressions
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type == token.LPAREN && p.PeekToken.Type == js.FUNCTION {
			node := &iifeExpr{LparenToken: p.CurrentToken}
			p.AdvanceToken()
			if node.Function, err = js.ParseFunctionDecl(p); err != nil {
				return
			}
			if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
				return
			}
			return node, nil
		}
		return next()
	})
	p := b.Build([]byte(input))
	result, err := js.ParseProgram(p)
	if err != nil {
		t.Fatal(err)
	}
	pr := xjs.NewPrinter()
	pr.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		if node, ok := node.(*iifeExpr); ok {
			p.Print(node.LparenToken)
			p.Print(node.Function)
			p.Print(node.RparenToken)
			return nil
		}
		return next(node)
	})
	pr.Init(printer.WithIndent("\t"))
	pr.Print(result)
	expected := "(\nfunction foo() {\n\tprint('Hello, World!');\n})();"
	require.Equal(t, expected, pr.String())
}
