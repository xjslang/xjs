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

func (p *Printer) Init(opts ...printerOption) {
	cfg := &printerConfig{
		indent: "  ",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if p.printer == nil {
		p.printer = p.defaultPrinter
	}
	p.indent = cfg.indent
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

func (p *Printer) Print(node ast.Node) {
	p.printer(node)
}

func (p *Printer) String() string {
	return p.doc.String()
}
