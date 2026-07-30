package internal

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/jsextended"
	"github.com/xjslang/xjs/printer"
)

func Debug(node ast.Node) (string, error) {
	pr := printer.NewBuilder().
		UsePrinter(js.Printer).
		UsePrinter(jsextended.Printer).
		UsePrinter(debugger).
		Build(printer.WithTrivia(false))
	pr.Print(node)
	return pr.Output()
}

// adds parentheses around expressions to make operator precedence explicit
// for debugging purposes
func debugger(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch node.(type) {
	case
		// core expressions to parenthesize in debug output
		*js.AssignExpr, *js.BinaryExpr, *js.CallExpr, *js.DecExpr, *js.IncExpr, *js.IndexExpr, *js.MemberExpr,
		*js.ArrayExpr, *js.FunctionExpr, *js.GroupExpr, *js.ObjExpr,
		// extended expressions to parenthesize in debug output
		*jsextended.TypeofExpr:
		pr.Print('(')
		defer pr.Print(')')
	}
	return next(node)
}
