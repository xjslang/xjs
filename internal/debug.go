package internal

import (
	"github.com/davecgh/go-spew/spew"
)

var cfg = &spew.ConfigState{
	Indent:                  "\t",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	ContinueOnMethod:        false,
}

// Dump prints the AST node structure for debugging
func Dump(node any) {
	cfg.Dump(node)
}
