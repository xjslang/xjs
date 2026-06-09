package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjslang/xjs/token"
)

func TestScanNumber(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
	}{
		{"integer", []string{"0", "123"}},
		{
			name:   "float",
			inputs: []string{"123.456", ".456", "123.", "123e32", "123.456e+34", ".456e-34"},
		},
		{
			name:   "hexadecimal",
			inputs: []string{"0x10", "0X20", "0xABCDEF", "0xabcdef", "0x123ABC", "0xFFFFFF", "0x0", "0xF"},
		},
		{
			name:   "octal",
			inputs: []string{"0o10", "0O20", "0o1234567", "0o0", "0o7", "0o777", "0o0012", "0O01234567"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, input := range test.inputs {
				sc := &Scanner{}
				sc.Init([]byte(input))
				result, typ := scanNumber(sc)
				if !assert.Equal(t, token.NUMBER, typ) {
					continue
				}
				assert.Equal(t, input, result)
			}
		})
	}

	t.Run("invalid formats", func(t *testing.T) {
		inputs := []string{
			"123e", "123e+", "123e-", // invalid float number
			"0x", // invalid hex number
			"0o", // invalid octal number
		}
		for _, input := range inputs {
			sc := &Scanner{}
			sc.Init([]byte(input))
			_, typ := scanNumber(sc)
			assert.Equal(t, token.ILLEGAL, typ)
		}
	})
}

func TestScanString(t *testing.T) {
	tests := []struct {
		delimiter rune
		input     string
	}{
		{'\'', `'string with \'escaped single quotes\''`},
		{'"', `"string with \"escaped quotes\""`},
	}
	for _, test := range tests {
		sc := &Scanner{}
		sc.Init([]byte(test.input))
		result, typ := scanString(sc, test.delimiter)
		if !assert.Equal(t, token.STRING, typ) {
			continue
		}
		assert.Equal(t, test.input, result)
	}
}
