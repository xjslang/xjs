package printer

import (
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
)

const eol = rune(-1) // end of line

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
	lastChar    rune
	ensureLine  bool
	ensureSpace bool
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
	p.lastChar = eol
}

func (p *Printer) PrintString(s string) {
	if len(s) == 0 {
		return
	}
	r, _ := utf8.DecodeLastRuneInString(s)
	p.lastChar = r
	p.doc.WriteString(s)
}

func (p *Printer) PrintRune(r rune) {
	p.lastChar = r
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

func (p *Printer) PrintNode(node ast.Node) {
	p.printer(node)
}

func (p *Printer) String() string {
	return p.doc.String()
}

func (p *Printer) PrintIndentedString(s string) {
	if len(s) == 0 {
		return
	}
	if p.ensureLine {
		p.printLineIfNeeded()
		p.ensureLine = false
	}
	if p.ensureSpace {
		p.printSpaceIfNeeded()
		p.ensureSpace = false
	}
	p.printIndentIfNeeded()
	p.PrintString(s)
}

// TODO: Comments on the "EOF" token are not being printed (will be addressed in another ticket)
// See: https://github.com/xjslang/xjs/pull/107#discussion_r3255438378
func (p *Printer) PrintTrivia(trivia []scanner.Token) {
	for _, tok := range trivia {
		switch tok.Type {
		case scanner.NEWLINE:
			p.PrintRune('\n')
		case scanner.LINE_COMMENT:
			p.printSpaceIfNeeded()
			p.printIndentIfNeeded()
			p.PrintString("//" + tok.Literal)
		case scanner.BLOCK_COMMENT:
			p.printIndentIfNeeded()
			p.PrintString("/*" + tok.Literal + "*/")
		}
	}
}

func (p *Printer) PrintToken(tok scanner.Token) {
	p.PrintTrivia(tok.LeadingTrivia)
	p.PrintIndentedString(tok.Literal)
}

func (p *Printer) EnsureLine() {
	p.ensureLine = true
}

func (p *Printer) EnsureSpace() {
	p.ensureSpace = true
}

func (p *Printer) printLineIfNeeded() {
	if !isNewLine(p.lastChar) {
		p.PrintRune('\n')
	}
}

func (p *Printer) printSpaceIfNeeded() {
	if !isWhitespace(p.lastChar) {
		p.PrintRune(' ')
	}
}

func (p *Printer) printIndentIfNeeded() {
	if isNewLine(p.lastChar) {
		p.PrintIndent()
	}
}
