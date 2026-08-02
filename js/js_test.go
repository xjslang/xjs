package js_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/internal"
	"github.com/xjslang/xjs/parser"
	"github.com/xorcare/golden"
)

func TestRoundtripFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/roundtrip/*.js")
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			// data -> AST
			result, err := xjs.Parse(data)
			require.NoError(t, err)

			// AST -> code
			code, err := xjs.Print(result)
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
			result, err := xjs.Parse(data)
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

			// data -> AST
			result, err := xjs.Parse(data)
			require.NoError(t, err)

			// AST -> code
			code, err := xjs.Print(result)
			require.NoError(t, err)

			// verify printed code parses
			_, err = xjs.Parse([]byte(code))
			require.NoError(t, err)
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
			_, errs := xjs.Parse(data)
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

func TestParseExpr_BinaryOpOnNextLine(t *testing.T) {
	input := "100\n+ 200"
	result, err := xjs.Parse([]byte(input))
	require.NoError(t, err)
	code, err := internal.Debug(result)
	require.NoError(t, err)
	require.Equal(t, "(100 + 200);", code)
}
