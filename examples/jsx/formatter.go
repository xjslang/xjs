package jsx

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
)

func Formatter(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *ConcatExpr:
		pr.Print(v.Left, " | ", v.Right)
	case *Tag:
		pr.SpPrint("<").Print(v.Name, ">")
		pr.IncreaseIndent()
		if v.Children != nil {
			pr.LnPrint(v.Children)
		}
		pr.DecreaseIndent()
		pr.LnPrint("</").Print(v.Name, ">")
	default:
		return next(node)
	}
	return nil
}
