package internal

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/xjslang/xjs/ast"
)

var cfg = &spew.ConfigState{
	Indent:                  "\t",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	ContinueOnMethod:        false,
}

// Dump prints the AST node structure for debugging
func Dump(node ast.Node) {
	cfg.Dump(node)
}
