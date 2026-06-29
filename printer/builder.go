package printer

import "github.com/xjslang/xjs/ast"

type Builder struct {
	printers []func(*Printer, ast.Node, func(ast.Node) error) error
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) UsePrinter(printer func(pr *Printer, node ast.Node, next func(node ast.Node) error) error) *Builder {
	b.printers = append(b.printers, printer)
	return b
}

func (b *Builder) Build(opts ...Option) *Printer {
	pr := &Printer{}
	for _, printer := range b.printers {
		pr.usePrinter(printer)
	}
	pr.init(opts...)
	return pr
}
