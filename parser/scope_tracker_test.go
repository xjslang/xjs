package parser

import (
	"testing"
)

func TestScopeTracker(t *testing.T) {
	st := scopeTracker{}
	if st.In(BlockScope) {
		t.Errorf("Expected not to be in blockScope")
	}
	for range 2 {
		st.Enter(BlockScope)
		if !st.In(BlockScope) {
			t.Errorf("Expected to be in blockScope")
		}
		// exit twice to verify that counter is not negative
		st.Exit(BlockScope)
		st.Exit(BlockScope)
	}
}
