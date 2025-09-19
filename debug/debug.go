package debug

import (
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/xjslang/xjs/ast"
)

var cfg = &spew.ConfigState{
	Indent:                  "\t",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	ContinueOnMethod:        false,
}

// ToString returns the textual representation of a node.
func ToString(node ast.Node) string {
	var b strings.Builder
	node.WriteTo(&b)
	return b.String()
}

// Print prints a formatted representation of a node.
func Print(node ast.Node) {
	cfg.Dump(node)
}
