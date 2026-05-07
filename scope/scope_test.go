package scope_test

import (
	"testing"

	"github.com/xjslang/xjs/scope"
)

func TestRegisterMultipleScopes(t *testing.T) {
	tr := scope.NewTracker()
	sc1 := tr.RegisterScope()
	sc2 := tr.RegisterScope()
	if sc1 == sc2 {
		t.Errorf("Expected scopes to be different")
	}
}

func TestScope(t *testing.T) {
	tr := scope.NewTracker()
	sc := tr.RegisterScope()
	if tr.In(sc) {
		t.Errorf("Expected not to be in blockScope")
	}
	for range 2 {
		tr.Enter(sc)
		if !tr.In(sc) {
			t.Errorf("Expected to be in blockScope")
		}
		tr.Exit(sc)
		tr.Exit(sc)
	}
}
