package scanner_test

import (
	"sync"
	"testing"

	"github.com/xjslang/xjs/scanner"
)

func TestConcurrentKindAccess(t *testing.T) {
	n := 100
	types := make([]scanner.Kind, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			types[i] = scanner.RegisterKind("aaa")
		}(i)
	}
	wg.Wait()

	// check that types has no duplicates
	seen := make(map[scanner.Kind]bool)
	for _, typ := range types {
		if seen[typ] {
			t.Fatalf("Duplicate!")
		}
		seen[typ] = true
	}
}
