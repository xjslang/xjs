package jsx

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
)

// Compiler transforms the code to valid JS code.
func Compiler(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *ConcatExpr:
		pr.Print("(function(){")
		pr.Print("const elem = document.createDocumentFragment();")
		pr.Print("elem.append(", v.Left, ");")
		pr.Print("elem.append(", v.Right, ");")
		pr.Print("return elem})()")
	case *Tag:
		pr.Print("(function(){")
		pr.Print("const elem = document.createElement('", v.Name, "');")
		if v.Children != nil {
			pr.Print("elem.append(", v.Children, ");")
		}
		pr.Print("return elem})()")
	default:
		return next(node)
	}
	return nil
}
