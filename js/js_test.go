package js_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
)

func TestRoundtrip(t *testing.T) {
	files, err := filepath.Glob("testdata/*.roundtrip.js")
	require.NoError(t, err)
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
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
