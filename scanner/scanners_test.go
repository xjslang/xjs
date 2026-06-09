package scanner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjslang/xjs/token"
)

func TestConsumeNumber(t *testing.T) {
	suffixes := []string{"", "e34", "e+34", "e-34"}
	tests := []string{"123", "123.456", ".456", "123."}
	for _, suffix := range suffixes {
		for _, test := range tests {
			input := fmt.Sprintf("%s%s", test, suffix)
			sc := &Scanner{}
			sc.Init([]byte(input))
			result, typ := scanNumber(sc)
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
			_, typ := scanNumber(sc)
			assert.Equal(t, token.ILLEGAL, typ)
		}
	})
}
