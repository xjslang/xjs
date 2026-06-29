package printer

import (
	"github.com/xjslang/xjs/ast"
)

func (pr *Printer) usePrinter(printer func(pr *Printer, node ast.Node, next func(node ast.Node) error) error) {
	print := pr.printer
	if pr.printer == nil {
		print = defaultPrinter
	}
	pr.printer = func(p *Printer, node ast.Node) error {
		return printer(p, node, func(node ast.Node) error {
			return print(p, node)
		})
	}
}

func defaultPrinter(p *Printer, node ast.Node) error {
	p.printString("<unknown>")
	return nil
}
