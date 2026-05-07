package parser

import (
	"testing"
)

func TestRegisterScope(t *testing.T) {
	sc1 := RegisterScope()
	sc2 := RegisterScope()
	if sc1 == sc2 {
		t.Errorf("Expected different scopes")
	}
}

func TestScopeTracker(t *testing.T) {
	sc := RegisterScope()
	st := scopeTracker{}
	if st.In(sc) {
		t.Errorf("Expected not to be in blockScope")
	}
	for range 2 {
		st.Enter(sc)
		if !st.In(sc) {
			t.Errorf("Expected to be in blockScope")
		}
		// exit twice to verify that counter is not negative
		st.Exit(sc)
		st.Exit(sc)
	}
}
