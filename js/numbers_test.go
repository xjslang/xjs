package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func TestScanHexNumber(t *testing.T) {
	tests := []string{"x10", "X20", "xABCDEF", "xabcdef", "x123ABC", "xFFFFFF", "x0", "xF"}
	for _, test := range tests {
		sc := &scanner.Scanner{}
		sc.Init([]byte(test))
		result, typ := ScanHexNumber(sc)
		if !assert.Equal(t, token.NUMBER, typ) {
			continue
		}
		assert.Equal(t, test, result)
	}

	t.Run("invalid formats", func(t *testing.T) {
		input := "x"
		sc := &scanner.Scanner{}
		sc.Init([]byte(input))
		_, typ := ScanHexNumber(sc)
		assert.Equal(t, token.ILLEGAL, typ)
	})
}
