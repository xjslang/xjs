package printer

import (
	"strings"
)

type Printer struct {
	doc strings.Builder
}

func (p *Printer) PrintString(s string) {
	p.doc.WriteString(s)
}

func (p *Printer) PrintRune(r rune) {
	p.doc.WriteRune(r)
}

func (p *Printer) String() string {
	return p.doc.String()
}
