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
	lastChar    rune
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
	p.lastChar = eol
}

func (p *Printer) PrintString(s string) {
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

func (p *Printer) Print(node ast.Node) {
	p.printer(node)
}

func (p *Printer) String() string {
	return p.doc.String()
}

// TODO: Comments on the "EOF" token are not being printed (will be addressed in another ticket)
func (p *Printer) PrintTrivia(trivia []scanner.Token) {
	for _, tok := range trivia {
		switch tok.Type {
		case scanner.NEWLINE:
			p.PrintRune('\n')
		case scanner.LINE_COMMENT:
			p.EnsureSpace()
			p.EnsureIndent()
			p.PrintString("//")
			p.PrintString(tok.Literal)
		case scanner.BLOCK_COMMENT:
			p.EnsureSpace()
			p.EnsureIndent()
			p.PrintString("/*")
			p.PrintString(tok.Literal)
			p.PrintString("*/")
		}
	}
}

// PrintTokenAfterNewline prints a token on a newline, ensuring a newline is added if needed.
// Leading trivia are printed before the newline.
func (p *Printer) PrintTokenAfterNewline(tok scanner.Token) {
	p.PrintTrivia(tok.LeadingTrivia)
	p.EnsureLine()
	p.EnsureIndent()
	p.PrintString(tok.Literal)
}

// PrintTokenAfterSpace prints a token with a leading space, ensuring a space is added if needed.
// Leading trivia are printed before the space.
func (p *Printer) PrintTokenAfterSpace(tok scanner.Token) {
	p.PrintTrivia(tok.LeadingTrivia)
	p.EnsureSpace()
	p.EnsureIndent()
	p.PrintString(tok.Literal)
}

// PrintToken prints a token at the current position with its leading trivia.
func (p *Printer) PrintToken(tok scanner.Token) {
	p.PrintTrivia(tok.LeadingTrivia)
	p.EnsureIndent()
	p.PrintString(tok.Literal)
}

// EnsureSpace prints a space only if the last character is not a space.
func (p *Printer) EnsureSpace() {
	// TODO: In the context of an expression, parentheses on the left have the same consideration as spaces. For example, "(100)" should not be formatted as "( 100)" (will be addressed in another ticket)
	if !isNewLine(p.lastChar) && p.lastChar != ' ' {
		p.PrintRune(' ')
	}
}

// EnsureLine prints a line only if the last character is not a newline.
func (p *Printer) EnsureLine() {
	if !isNewLine(p.lastChar) {
		p.PrintRune('\n')
	}
}

// EnsureIndent prints an indentation only if the last character is a newline.
func (p *Printer) EnsureIndent() {
	if isNewLine(p.lastChar) {
		p.PrintIndent()
	}
}
