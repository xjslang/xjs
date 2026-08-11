package jsextended_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/internal"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/jsextended"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xorcare/golden"
)

func Parse(input []byte) (*js.Program, error) {
	sb := jsextended.ScannerBuilder()
	sc := sb.Build(input)
	pb := jsextended.ParserBuilder()
	p := pb.Build(sc)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := jsextended.PrinterBuilder().Build(opts...)
	pr.Print(result)
	return pr.Output()
}

func BenchmarkParse(b *testing.B) {
	dat, err := os.ReadFile("testdata/check/general.js")
	require.NoError(b, err)
	require.NotEmpty(b, dat)

	// verify that it can be parsed and printed
	result, err := Parse(dat)
	require.NoError(b, err)
	_, err = Print(result)
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		_, _ = Parse(dat)
	}
}

func BenchmarkPrint(b *testing.B) {
	dat, err := os.ReadFile("testdata/check/general.js")
	require.NoError(b, err)
	require.NotEmpty(b, dat)

	// verify that it can be parsed and printed
	result, err := Parse(dat)
	require.NoError(b, err)
	_, err = Print(result)
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		_, _ = Print(result)
	}
}

func BenchmarkParseAndPrint(b *testing.B) {
	dat, err := os.ReadFile("testdata/check/general.js")
	require.NoError(b, err)
	require.NotEmpty(b, dat)

	// verify that it can be parsed and printed
	result, err := Parse(dat)
	require.NoError(b, err)
	_, err = Print(result)
	require.NoError(b, err)

	b.ResetTimer()
	for b.Loop() {
		var err error
		result, err = Parse(dat)
		if err != nil {
			b.Fatal(err)
		}
		_, err = Print(result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestRoundtripFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/roundtrip/*.js")
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			// data -> AST
			result, err := Parse(data)
			require.NoError(t, err)

			// AST -> code
			code, err := Print(result)
			require.NoError(t, err)

			assert.Equal(t, string(data), code)
		})
	}
}

func TestDebugFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/debug/*.js")
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			// data -> AST
			result, err := Parse(data)
			require.NoError(t, err)

			// AST -> code
			code, err := internal.Debug(result)
			require.NoError(t, err)

			golden.Assert(t, []byte(code))
		})
	}
}

func TestCheckFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/check/*.js")
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			result1, err := Parse(data)
			require.NoError(t, err)

			code1, err := Print(result1, printer.Compact())
			require.NoError(t, err)

			result2, err := Parse([]byte(code1))
			require.NoError(t, err)

			code2, err := Print(result2, printer.Compact())
			require.NoError(t, err)
			require.Equal(t, code1, code2)
		})
	}
}

func TestErrorFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/errors/*.js")
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			// data -> AST
			_, errs := Parse(data)
			require.Error(t, errs)
			require.IsType(t, parser.ErrorList{}, errs)

			var msgs []string
			errList := errs.(parser.ErrorList)
			for _, err := range errList {
				require.IsType(t, parser.Error{}, err)
				e := err.(parser.Error)
				msgs = append(msgs, e.Message)
			}

			golden.Assert(t, []byte(strings.Join(msgs, "\n")))
		})
	}
}
