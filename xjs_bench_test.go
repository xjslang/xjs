package xjs_test

import (
	"os"
	"testing"

	"github.com/dop251/goja/parser"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/parse/v2"
	tdewolffjs "github.com/tdewolff/parse/v2/js"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

// used to prevent the compiler from optimizing the benchmark code
var sink any

func Parse(input string) (*js.Program, error) {
	sb := js.ScannerBuilder()
	sc := sb.Build(input)
	pb := js.ParserBuilder()
	p := pb.Build(sc)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := js.PrinterBuilder().Build(opts...)
	pr.Print(result)
	return pr.Output()
}

func BenchmarkCompare(b *testing.B) {
	dat, err := os.ReadFile("testdata/bench/lodash.js")
	require.NoError(b, err)
	require.NotEmpty(b, dat)

	b.Run("parser=xjs", func(b *testing.B) {
		// verify that it can be parsed
		content := string(dat)
		sink, err = Parse(content)
		require.NoError(b, err)

		b.SetBytes(int64(len(dat)))
		b.ResetTimer()
		for b.Loop() {
			sink, err = Parse(content)
		}
	})

	b.Run("parser=goja", func(b *testing.B) {
		// verify that the file can be parsed
		content := string(dat)
		sink, err = parser.ParseFile(nil, "", content, 0)
		require.NoError(b, err)

		b.SetBytes(int64(len(dat)))
		b.ResetTimer()
		for b.Loop() {
			sink, err = parser.ParseFile(nil, "", content, 0)
		}
	})

	b.Run("parser=tdewolff", func(b *testing.B) {
		// verify that the file can be parsed
		content := string(dat)
		sink, err = tdewolffjs.Parse(parse.NewInputString(content), tdewolffjs.Options{})
		require.NoError(b, err)

		b.SetBytes(int64(len(dat)))
		b.ResetTimer()
		for b.Loop() {
			sink, err = tdewolffjs.Parse(parse.NewInputString(content), tdewolffjs.Options{})
		}
	})
}
