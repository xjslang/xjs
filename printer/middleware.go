package printer

import (
	"github.com/xjslang/xjs/ast"
)

func (p *Printer) UsePrinter(printer func(p *Printer, node ast.Node, next func(node ast.Node))) {
	print := p.printer
	if p.printer == nil {
		print = defaultPrinter
	}
	p.printer = func(p *Printer, node ast.Node) {
		printer(p, node, func(node ast.Node) {
			print(p, node)
		})
	}
}

func defaultPrinter(p *Printer, node ast.Node) {
	p.printString("<unknown>")
}
