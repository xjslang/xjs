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
	"github.com/xorcare/golden"
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
	out, err := xjs.Print(result)
	if err != nil {
		panic(err)
	}
	fmt.Print(out)
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}

func TestParserErrors(t *testing.T) {
	input := `// program stmt
aaa bbb // ; or newline expected

// assign stmt
x =; // expression expected

// block stmt
{ aaa bbb } // ; or newline expected

// for stmt
for; // ( expected
for (?; // init expected
for (let i = 0; i < 10; i++; // ) expected

// function stmt
function; // identifier expected
function foo; // ( expected
function foo(; // identifier expected
function foo(a,; // identifier expected
function foo(p0; // ) expected

// if stmt
if; // ( expected
if (; // expression expected

// let stmt
let; // identifier expected
let x; // = expected
let x =; // expression expected

// arr expr
[; // expression expected
[1,; // expression expected
[1; // ] expected

// binary expr
1+; // expression expected

// unary expr
!; // expression expected

// call expr
print(; // expression expected
print(1, // expression expected
print(1; // ) expected

// function expr
(function; // ( expected
(function(; // identifier expected
(function(a,; // identifier expected
(function(a; // ) expected

// group expr
(100; // ) expected

// index expr
a[; // expression expected
a[100; // ] expected

// obj expr
({; // key expected
({name; // : expected
({name:; // expression expected
({name: 100,; // key expected
({name: 100; // } expected

// numbers
.123; // expression expected (numbers cannot start with '.')
1x123; // ; or newline expected (invalid hex)
2O123; // ; or newline expected (invalid octal)
0X; // expression expected (incomplete hex)
0o; // expression expected (incomplete octal)

// member expr
a.100; // key expected
a.(b); // key expected
a.(b + c); // key expected

// reserved keys cannot be used as identifiers
let if = 100; // identifier expected`
	_, errs := xjs.Parse([]byte(input))
	require.IsType(t, parser.ErrorList{}, errs)
	golden.Assert(t, []byte(errs.Error()))

	t.Run("key expected at eof", func(t *testing.T) {
		_, err := xjs.Parse([]byte("a."))
		require.EqualError(t, err, "[line:0, col:2] key expected")
	})
}

func TestLanguageFeatures(t *testing.T) {
	pattern := filepath.Join("testdata", "*.js")
	files, err := filepath.Glob(pattern)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		testName := strings.TrimSuffix(filepath.Base(file), ".js")
		t.Run(testName, func(t *testing.T) {
			// read file
			dat, err := os.ReadFile(file)
			require.NoError(t, err)
			// parse data
			result, err := xjs.Parse(dat)
			require.NoError(t, err)
			// print result
			out, err := xjs.Print(result)
			require.NoError(t, err)
			// re-parse the output
			result, err = xjs.Parse([]byte(out))
			require.NoError(t, err)
			// re-print the result
			out, err = xjs.Print(result)
			require.NoError(t, err)
			// the original must match the final printed result
			require.Equal(t, string(dat), out)
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
			expected: "let a = point(\nx,\ny\n);",
		},
		{
			input:    "let a = setTimeout(function() { console.log('tick!') }, 1000,)",
			expected: "let a = setTimeout(function () {\n\tconsole.log('tick!');\n}, 1000);",
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
			out, err := xjs.Print(result, printer.WithIndent("\t"))
			require.NoError(t, err)
			assert.Equal(t, test.expected, out)
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
	b := xjs.PluginBuilder()
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
	pr := xjs.PrinterBuilder().
		UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
			if node, ok := node.(*iifeExpr); ok {
				p.Print(node.LparenToken)
				p.Print(node.Function)
				p.Print(node.RparenToken)
				return nil
			}
			return next(node)
		}).
		Build(printer.WithIndent("\t"))
	pr.Print(result)
	expected := "(\nfunction foo() {\n\tprint('Hello, World!');\n})();"
	out, err := pr.Output()
	require.NoError(t, err)
	require.Equal(t, expected, out)
}
