package printer

import (
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

const eol = rune(-1)

type config struct {
	indent string
}

type Option func(*config)

func WithIndent(value string) Option {
	return func(cfg *config) {
		cfg.indent = value
	}
}

type Printer struct {
	doc         strings.Builder
	indent      string
	indentLevel int
	lastChar    rune
	ensureLine  bool
	ensureSpace bool
	printer     func(*Printer, ast.Node)
}

func (p *Printer) Init(opts ...Option) {
	cfg := &config{
		indent: "  ",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	p.doc.Reset()
	p.indent = cfg.indent
	p.indentLevel = 0
	p.lastChar = eol
	if p.printer == nil {
		p.printer = defaultPrinter
	}
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
		p.writeString(p.indent)
	}
}

func (p *Printer) EnsureLine() {
	p.ensureLine = true
}

func (p *Printer) EnsureSpace() {
	p.ensureSpace = true
}

func (p *Printer) Print(args ...any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			p.printString(v)
		case rune:
			p.printRune(v)
		case ast.Node:
			p.printNode(v)
		case token.Token:
			p.printToken(v)
		default:
			panic("Unsupported type")
		}
	}
}

func (p *Printer) LnPrint(arg any) {
	p.EnsureLine()
	p.Print(arg)
}

func (p *Printer) SpPrint(args ...any) {
	for _, arg := range args {
		p.EnsureSpace()
		p.Print(arg)
	}
}

func (p *Printer) PrintTrivia(trivia []token.Token) {
	for _, tok := range trivia {
		switch tok.Type {
		case token.NEWLINE:
			p.writeRune('\n')
		case token.LINE_COMMENT:
			p.printSpaceIfNeeded()
			p.printIndentIfNeeded()
			p.writeString("//" + tok.Literal)
		case token.BLOCK_COMMENT:
			p.printIndentIfNeeded()
			p.writeString("/*" + tok.Literal + "*/")
		}
	}
}

func (p *Printer) String() string {
	return p.doc.String()
}

func (p *Printer) Bytes() []byte {
	return []byte(p.String())
}

func (p *Printer) writeString(s string) {
	if len(s) == 0 {
		return
	}
	r, _ := utf8.DecodeLastRuneInString(s)
	p.lastChar = r
	p.doc.WriteString(s)
}

func (p *Printer) writeRune(r rune) {
	p.lastChar = r
	p.doc.WriteRune(r)
}

func (p *Printer) printNode(node ast.Node) {
	p.printer(p, node)
}

func (p *Printer) printString(s string) {
	if len(s) == 0 {
		return
	}
	p.printSeparatorIfNeeded()
	p.writeString(s)
}

func (p *Printer) printRune(r rune) {
	p.printSeparatorIfNeeded()
	p.writeRune(r)
}

func (p *Printer) printToken(tok token.Token) {
	p.PrintTrivia(tok.LeadingTrivia)
	p.printString(tok.Literal)
}

func (p *Printer) printLineIfNeeded() {
	if !isNewLine(p.lastChar) {
		p.writeRune('\n')
	}
}

func (p *Printer) printSpaceIfNeeded() {
	if !isWhitespace(p.lastChar) {
		p.writeRune(' ')
	}
}

func (p *Printer) printIndentIfNeeded() {
	if isNewLine(p.lastChar) {
		p.PrintIndent()
	}
}

func (p *Printer) printSeparatorIfNeeded() {
	if p.ensureLine {
		p.printLineIfNeeded()
		p.ensureLine = false
	}
	if p.ensureSpace {
		p.printSpaceIfNeeded()
		p.ensureSpace = false
	}
	p.printIndentIfNeeded()
}
