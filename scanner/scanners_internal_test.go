package scanner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/scanner"
)

func TestScanString(t *testing.T) {
	tests := []struct {
		delimiter rune
		input     string
	}{
		{'\'', `'string with \'escaped single quotes\''`},
		{'"', `"string with \"escaped quotes\""`},
	}
	for _, test := range tests {
		sc := &scanner.Scanner{}
		sc.Init([]byte(test.input))
		sc.AdvanceChar()
		result, err := js.ScanString(sc, test.delimiter)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, test.input, string(test.delimiter)+result)
	}
}
