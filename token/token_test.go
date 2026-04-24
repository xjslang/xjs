package token

import (
	"sync"
	"testing"
)

func TestRegisterType(t *testing.T) {
	lit := "**"
	powType := RegisterType(lit)
	if powType < initCustomType {
		t.Errorf("Expected type to be greater than %d, got %d", initCustomType, powType)
	}
	if powType.String() != lit {
		t.Errorf("Expected %q, got %q", lit, powType.String())
	}
}

func TestRegisterType_ConcurrentReturnUniqueKeys(t *testing.T) {
	n := 100
	types := make([]TokenType, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			types[i] = RegisterType("aaa")
		}(i)
	}
	wg.Wait()

	// check that types has no duplicates
	seen := make(map[TokenType]bool)
	for _, typ := range types {
		if seen[typ] {
			t.Fatalf("Duplicate!")
		}
		seen[typ] = true
	}
}
