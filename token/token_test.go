package token_test

import (
	"sync"
	"testing"

	"github.com/xjslang/xjs/token"
)

func TestConcurrentKindAccess(t *testing.T) {
	n := 100
	types := make([]token.Type, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			types[i] = token.RegisterKind("aaa")
		}(i)
	}
	wg.Wait()

	// check that types has no duplicates
	seen := make(map[token.Type]bool)
	for _, typ := range types {
		if seen[typ] {
			t.Fatalf("Duplicate!")
		}
		seen[typ] = true
	}
}
