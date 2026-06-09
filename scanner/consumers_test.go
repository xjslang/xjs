package scanner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjslang/xjs/token"
)

func TestConsumeDecimalNumber(t *testing.T) {
	suffixes := []string{"", "e34", "e+34", "e-34"}
	tests := []string{"123", "123.456", ".456", "123."}
	for _, suffix := range suffixes {
		for _, test := range tests {
			input := fmt.Sprintf("%s%s", test, suffix)
			sc := &Scanner{}
			sc.Init([]byte(input))
			result, typ := sc.consumeDecimalNumber()
			if !assert.Equal(t, token.NUMBER, typ) {
				continue
			}
			assert.Equal(t, input, result)
		}
	}

	t.Run("invalid number formats", func(t *testing.T) {
		tests := []string{"123e", "123e+", "123e-"}
		for _, test := range tests {
			sc := &Scanner{}
			sc.Init([]byte(test))
			_, typ := sc.consumeDecimalNumber()
			assert.Equal(t, token.ILLEGAL, typ)
		}
	})
}

func TestConsumeHexNumber(t *testing.T) {
	tests := []string{"0x10", "0X20", "0xABCDEF", "0xabcdef", "0x123ABC", "0xFFFFFF", "0x0", "0xF"}
	for _, test := range tests {
		sc := &Scanner{}
		sc.Init([]byte(test))
		result, typ := sc.consumeNumber()
		if !assert.Equal(t, token.NUMBER, typ) {
			continue
		}
		assert.Equal(t, test, result)
	}

	t.Run("invalid formats", func(t *testing.T) {
		input := "0x"
		sc := &Scanner{}
		sc.Init([]byte(input))
		_, typ := sc.consumeNumber()
		assert.Equal(t, token.ILLEGAL, typ)
	})
}
