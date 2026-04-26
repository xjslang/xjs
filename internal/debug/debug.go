package debug

import (
	"github.com/davecgh/go-spew/spew"
)

var cfg = &spew.ConfigState{
	Indent:                  "   ",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	ContinueOnMethod:        false,
}

func Print(node ...any) {
	cfg.Dump(node)
}
