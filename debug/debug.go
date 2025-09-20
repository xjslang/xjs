package debug

import (
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/xjslang/xjs/ast"
)

var cfg = &spew.ConfigState{
	Indent:                  "   ",
	DisableMethods:          true,
	DisablePointerAddresses: true,
	ContinueOnMethod:        false,
}

// ToString converts an AST node to its string representation.
func ToString(node ast.Node) string {
	var b strings.Builder
	node.WriteTo(&b)
	return b.String()
}

// Print outputs a detailed formatted representation of an AST node for debugging.
func Print(node ast.Node) {
	cfg.Dump(node)
}
