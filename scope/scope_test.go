package scope_test

import (
	"testing"

	"github.com/xjslang/xjs/scope"
)

func TestRegisterMultipleScopes(t *testing.T) {
	sc1 := scope.RegisterScope()
	sc2 := scope.RegisterScope()
	if sc1 == sc2 {
		t.Errorf("Expected scopes to be different")
	}
}

func TestScope(t *testing.T) {
	sc := scope.RegisterScope()
	if scope.In(sc) {
		t.Errorf("Expected not to be in blockScope")
	}
	for range 2 {
		scope.Enter(sc)
		if !scope.In(sc) {
			t.Errorf("Expected to be in blockScope")
		}
		scope.Exit(sc)
		scope.Exit(sc)
	}
}
