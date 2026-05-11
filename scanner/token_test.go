package scanner

import (
	"sync"
	"testing"
)

func TestRegisterKind(t *testing.T) {
	lit := "**"
	powType := RegisterKind(lit)
	if powType < initCustomType {
		t.Errorf("Expected type to be greater than %d, got %d", initCustomType, powType)
	}
	if powType.String() != lit {
		t.Errorf("Expected %q, got %q", lit, powType.String())
	}

	t.Run("concurrent access returns unique keys", func(t *testing.T) {
		n := 100
		types := make([]Kind, n)
		var wg sync.WaitGroup
		for i := range n {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				types[i] = RegisterKind("aaa")
			}(i)
		}
		wg.Wait()

		// check that types has no duplicates
		seen := make(map[Kind]bool)
		for _, typ := range types {
			if seen[typ] {
				t.Fatalf("Duplicate!")
			}
			seen[typ] = true
		}
	})
}
