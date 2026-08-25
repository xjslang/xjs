package internal

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

func Debug(node ast.Node) (string, error) {
	prb := js.PrinterBuilder()
	prb.UsePrinter(debugger)
	pr := prb.Build(printer.Compact())
	pr.Print(node)
	return pr.Output()
}

// adds parentheses around expressions to make operator precedence explicit
// for debugging purposes
func debugger(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
	switch node.(type) {
	case
		*js.InfixAssignOp, *js.InfixOp, *js.InfixParenOp, *js.PostfixDecOp, *js.PostfixIncOp, *js.InfixBracketOp, *js.InfixDotOp,
		*js.PrefixBracketOp, *js.PrefixFunctionOp, *js.PrefixParenOp, *js.PrefixBraceOp,
		*js.PrefixTypeofOp:
		pr.Print('(')
		defer pr.Print(')')
	}
	return next(pr, node)
}
