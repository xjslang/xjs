package debug

import (
	"github.com/davecgh/go-spew/spew"
)

var cfg = spew.ConfigState{
	Indent:                  "\t",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	ContinueOnMethod:        false,
}

func Print(node ...any) {
	cfg.Dump(node...)
}

func Sprint(node ...any) string {
	return cfg.Sdump(node...)
}
