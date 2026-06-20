package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		sc := &Scanner{}
		sc.Init([]byte(test.input))
		result, err := scanString(sc, test.delimiter)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, test.input, result)
	}
}
