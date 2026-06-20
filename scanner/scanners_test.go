package scanner_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/scanner"
)

func TestScanNumber(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
	}{
		{"integer", []string{"0", "123"}},
		{
			name:   "float",
			inputs: []string{"123.456", ".456", "123.", "123e32", "123.456e+34", ".456e-34", "1e2"},
		},
		{
			name:   "hexadecimal",
			inputs: []string{"x10", "X20", "xABCDEF", "xabcdef", "x123ABC", "xFFFFFF", "x0", "xF"},
		},
		{
			name:   "octal",
			inputs: []string{"o10", "O20", "o1234567", "o0", "o7", "o777", "o0012", "O01234567"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, input := range test.inputs {
				sc := &scanner.Scanner{}
				sc.Init([]byte(input))
				var result string
				var err error
				switch sc.CurrentChar() {
				case 'x', 'X':
					result, err = js.ScanHexNumber(sc)
				case 'o', 'O':
					result, err = js.ScanOctalNumber(sc)
				default:
					result, err = js.ScanNumber(sc)
				}
				assert.NoError(t, err)
				assert.Equal(t, input, result)
			}
		})
	}

	t.Run("invalid formats", func(t *testing.T) {
		inputs := []string{
			"123e", "123e+", "123e-", "1e", // invalid float numbers
			"x", // invalid hex number
			"o", // invalid octal number
		}
		for _, input := range inputs {
			sc := &scanner.Scanner{}
			sc.Init([]byte(input))
			var err error
			switch sc.CurrentChar() {
			case 'x', 'X':
				_, err = js.ScanHexNumber(sc)
			case 'o', 'O':
				_, err = js.ScanOctalNumber(sc)
			default:
				_, err = js.ScanNumber(sc)
			}
			assert.Error(t, err)
		}
	})
}
