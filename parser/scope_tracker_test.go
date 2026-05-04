package parser

import (
	"testing"
)

func TestXxx(t *testing.T) {
	st := scopeTracker{}
	if st.in(blockScope) {
		t.Errorf("Expected not to be in blockScope")
	}
	for range 2 {
		st.enter(blockScope)
		if !st.in(blockScope) {
			t.Errorf("Expected to be in blockScope")
		}
		// exit twice to verify that counter is not negative
		st.exit(blockScope)
		st.exit(blockScope)
	}
}
