package xjs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/internal"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func TestExportStmt(t *testing.T) {
	t.Run("declaration expected", func(t *testing.T) {
		input := `export (function () { console.log('foo') })()`
		_, err := xjs.Parse([]byte(input))
		require.Error(t, err)
		require.Equal(t, err.Error(), "[line:0, col:7] declaration expected")
	})
}

func TestStandaloneSemicolons(t *testing.T) {
	input := `; // c1
	; // c2
	aaa()
	;
	bbb();
	ccc()`
	result, err := xjs.Parse([]byte(input))
	require.NoError(t, err)
	out, err := xjs.Print(result, printer.WithTrivia(false))
	require.NoError(t, err)
	require.Equal(t, `;
;
aaa();
bbb();
ccc();`, out)
}

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

func TestExpectSemi(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedToks []token.Token // current and next token
	}{
		{"consume semi", ";aaa", []token.Token{
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.IDENT, Literal: "aaa"},
		}},
		{"consume semi after newline", "\n;aaa", []token.Token{
			{Type: token.SEMICOLON, Literal: ";", AfterNewline: true},
			{Type: token.IDENT, Literal: "aaa"},
		}},
		{"do not consume }", "}aaa", []token.Token{
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.RBRACE, Literal: "}"},
		}},
		{"do not consume )", ")aaa", []token.Token{
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.RPAREN, Literal: ")"},
		}},
		{"return synthetic token on newline", "\naaa", []token.Token{
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.IDENT, Literal: "aaa", AfterNewline: true},
		}},
		{"return synthetic token on end of line", "", []token.Token{
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.EOF, Literal: ""},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := scanner.NewBuilder().Build([]byte(test.input))
			p := parser.NewBuilder().Build(s)
			tok, err := js.ExpectSemi(p)
			require.NoError(t, err)
			internal.AssertTokens(t, []token.Token{tok, p.CurrentToken}, test.expectedToks, internal.CompareAfterNewline())
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
		UsePrinter(func(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
			if node, ok := node.(*iifeExpr); ok {
				pr.Print(node.LparenToken)
				pr.Print(node.Function)
				pr.Print(node.RparenToken)
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
