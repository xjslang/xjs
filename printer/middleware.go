package printer

import (
	"github.com/xjslang/xjs/ast"
)

func (pr *Printer) usePrinter(printer func(pr *Printer, node ast.Node, next func(*Printer, ast.Node) error) error) {
	next := pr.printer
	if pr.printer == nil {
		next = defaultPrinter
	}
	pr.printer = func(p *Printer, node ast.Node) error {
		return printer(p, node, next)
	}
}

func defaultPrinter(p *Printer, node ast.Node) error {
	p.printString("<unknown>")
	return nil
}
