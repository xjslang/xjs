package internal

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/jsextended"
	"github.com/xjslang/xjs/printer"
)

func Debug(node ast.Node) (string, error) {
	prb := jsextended.PrinterBuilder()
	prb.UsePrinter(debugger)
	pr := prb.Build(printer.WithTrivia(false))
	pr.Print(node)
	return pr.Output()
}

// adds parentheses around expressions to make operator precedence explicit
// for debugging purposes
func debugger(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch node.(type) {
	case
		// core expressions to parenthesize in debug output
		*js.InfixAssignOp, *js.InfixOp, *js.InfixParenOp, *js.PostfixDecOp, *js.PostfixIncOp, *js.InfixBracketOp, *js.InfixDotOp,
		*js.PrefixBracketOp, *js.PrefixFunctionOp, *js.PrefixParenOp, *js.PrefixBraceOp,
		// extended expressions to parenthesize in debug output
		*jsextended.PrefixTypeofOp:
		pr.Print('(')
		defer pr.Print(')')
	}
	return next(node)
}
