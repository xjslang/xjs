package printer

import (
	"strings"

	"github.com/xjslang/xjs/ast"
)

type printerConfig struct {
	indent string
}

type printerOption func(*printerConfig)

func WithIndent(value string) printerOption {
	return func(cfg *printerConfig) {
		cfg.indent = value
	}
}

type Printer struct {
	doc         strings.Builder
	indent      string
	indentLevel int
	printer     func(ast.Node)
}

func New(opts ...printerOption) *Printer {
	cfg := &printerConfig{
		indent: "  ",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Printer{
		indent: cfg.indent,
	}
}

func (p *Printer) PrintString(s string) {
	p.doc.WriteString(s)
}

func (p *Printer) PrintRune(r rune) {
	p.doc.WriteRune(r)
}

func (p *Printer) IncreaseIndent() {
	p.indentLevel++
}

func (p *Printer) DecreaseIndent() {
	if p.indentLevel > 0 {
		p.indentLevel--
	}
}

func (p *Printer) PrintIndent() {
	for range p.indentLevel {
		p.doc.WriteString(p.indent)
	}
}

func (p *Printer) UsePrinter(printer func(c *Printer, node ast.Node, next func())) {
	print := p.printer
	if p.printer == nil {
		print = func(node ast.Node) {
			p.PrintString("<" + node.Type() + ">")
		}
	}
	p.printer = func(node ast.Node) {
		printer(p, node, func() {
			print(node)
		})
	}
}

func (p *Printer) Print(node ast.Node) {
	p.printer(node)
}

func (p *Printer) String() string {
	return p.doc.String()
}
