//go:build integration

package jsextended_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs/printer"
)

func TestPassFiles(t *testing.T) {
	files, err := filepath.Glob("testdata/pass/*.js")
	require.NoError(t, err)
	require.NotEmpty(t, files)
	passed, total := 0, 0
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			total++
			result1, err := Parse([]byte(data))
			require.NoError(t, err)

			code1, err := Print(result1, printer.Compact())
			require.NoError(t, err)

			result2, err := Parse([]byte(code1))
			require.NoError(t, err)

			code2, err := Print(result2, printer.Compact())
			require.NoError(t, err)
			if assert.Equal(t, code1, code2) {
				passed++
			}
		})
	}
	t.Logf("Total: %d, Passed: %d, Percent: %.2f", total, passed, 100*float64(passed)/float64(total))
}
