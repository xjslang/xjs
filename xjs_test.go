package xjs_test

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var updateGoldenFiles bool

func TestMain(m *testing.M) {
	flag.BoolVar(&updateGoldenFiles, "update", false, "update golden files")
	flag.Parse()
	os.Exit(m.Run())
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

func TestGoldenFiles(t *testing.T) {
	if updateGoldenFiles {
		t.Log("updating golden files")
	}
	files, err := filepath.Glob("./testdata/*.js")
	require.NoError(t, err)
	for _, file := range files {
		ext := filepath.Ext(file)
		goldFile := fmt.Sprintf("%s.print.txt", strings.TrimSuffix(file, ext))
		if !updateGoldenFiles && !assert.FileExists(t, goldFile) {
			continue
		}
		// parse the source file
		source, err := os.ReadFile(file)
		require.NoError(t, err)
		p := xjs.NewBuilder().Build(source)
		result, err := p.Parse()
		require.NoError(t, err)
		// print the result
		pr := xjs.NewPrinter()
		pr.Init()
		pr.Print(result)
		// create or update golden file
		got := pr.Bytes()
		if updateGoldenFiles {
			err = os.WriteFile(goldFile, got, 0o644)
			require.NoError(t, err)
			continue
		}
		// compare golden file with `got`
		want, err := os.ReadFile(goldFile)
		require.NoError(t, err)
		if got, want := string(got), string(want); got != want {
			diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
				A:       difflib.SplitLines(got),
				B:       difflib.SplitLines(want),
				Context: 5,
			})
			assert.NoError(t, err)
			t.Error(diff)
		}
	}
}

type iifeExpr struct {
	LparenToken token.Token
	RparenToken token.Token
	Function    *js.FunctionDecl
}

func (node *iifeExpr) Type() string {
	return "iifeExpr"
}

func TestMiddlewares(t *testing.T) {
	input := `(function foo() {
		print('Hello, World!')
	})()`
	b := xjs.NewBuilder()
	// parse IIFE expressions
	b.UsePrefixParser(func(p *parser.Parser, next func() (ast.Node, error)) (_ ast.Node, err error) {
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
	result, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	pr := xjs.NewPrinter()
	pr.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		if node, ok := node.(*iifeExpr); ok {
			p.Print(node.LparenToken)
			p.Print(node.Function)
			p.Print(node.RparenToken)
			return
		}
		next(node)
	})
	pr.Init(printer.WithIndent("\t"))
	pr.Print(result)
	expected := "(\nfunction foo() {\n\tprint('Hello, World!');\n})();"
	require.Equal(t, expected, pr.String())
}
