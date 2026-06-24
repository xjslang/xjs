package xjs_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/token"
)

func TestIfElseNonBlockThen(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "else after single-statement then",
			input:    `if (x) foo() else bar()`,
			expected: `if (x) foo() else bar();`,
		},
		{
			name:     "else after single-statement then with newlines",
			input:    "if (x)\n  foo()\nelse\n  bar()",
			expected: "if (x)\nfoo()\nelse\nbar();",
		},
		{
			name:     "else-if chain with non-block first then",
			input:    `if (a) foo() else if (b) bar() else baz()`,
			expected: `if (a) foo() else if (b) bar() else baz();`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := xjs.Parse([]byte(tt.input))
			require.NoError(t, err)
			pr := xjs.NewPrinter()
			pr.Init()
			pr.Print(result)
			out, err := pr.Output()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, out)
		})
	}
}

func TestDollarSignIdentifier(t *testing.T) {
	tests := []string{
		`let $el = 1`,
		`let $$ = 2`,
		`$scope.apply()`,
	}
	for _, src := range tests {
		t.Run(src, func(t *testing.T) {
			_, err := xjs.Parse([]byte(src))
			assert.NoError(t, err)
		})
	}
}

func TestUnicodeIdentifier(t *testing.T) {
	tests := []string{
		`let café = 1`,
		`let π = 3`,
		`let über = 4`,
	}
	for _, src := range tests {
		t.Run(src, func(t *testing.T) {
			_, err := xjs.Parse([]byte(src))
			assert.NoError(t, err)
		})
	}
}

func TestRegisterBinaryOpRace(t *testing.T) {
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			typ := token.RegisterBinaryOp("~~", 5)
			assert.True(t, typ.IsBinaryOp())
		}()
	}
	wg.Wait()
}

func TestRegisterUnaryOpRace(t *testing.T) {
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			typ := token.RegisterUnaryOp("^^")
			assert.True(t, typ.IsUnaryOp())
		}()
	}
	wg.Wait()
}

func TestScannerErrExposed(t *testing.T) {
	_, err := xjs.Parse([]byte(`"unterminated`))
	require.Error(t, err)
}

func TestBuilderBuildIndependence(t *testing.T) {
	b := xjs.NewBuilder()

	p1 := b.Build([]byte(`let x = 1`))
	p2 := b.Build([]byte(`let y = 2`))

	prog1, err := js.ParseProgram(p1)
	require.NoError(t, err)

	prog2, err := js.ParseProgram(p2)
	require.NoError(t, err)

	pr1 := xjs.NewPrinter()
	pr1.Init()
	pr1.Print(prog1)
	out1, _ := pr1.Output()

	pr2 := xjs.NewPrinter()
	pr2.Init()
	pr2.Print(prog2)
	out2, _ := pr2.Output()

	assert.True(t, strings.Contains(out1, "x"))
	assert.True(t, strings.Contains(out2, "y"))
	assert.False(t, strings.Contains(out1, "y"))
}
