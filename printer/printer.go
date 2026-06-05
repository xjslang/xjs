package printer

import (
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

const eol = rune(-1)

type config struct {
	indent        string
	lineComments  bool
	blockComments bool
	newLines      bool
}

type Option func(*config)

func Compact() Option {
	return func(cfg *config) {
		cfg.lineComments = false
		cfg.blockComments = false
		cfg.newLines = false
	}
}

func WithIndent(value string) Option {
	return func(cfg *config) {
		cfg.indent = value
	}
}

func WithComments(value bool) Option {
	return func(cfg *config) {
		cfg.lineComments = value
		cfg.blockComments = value
	}
}

func WithLineComments(value bool) Option {
	return func(cfg *config) {
		cfg.lineComments = value
	}
}

func WithBlockComments(value bool) Option {
	return func(cfg *config) {
		cfg.blockComments = value
	}
}

func WithNewLines(value bool) Option {
	return func(cfg *config) {
		cfg.newLines = value
	}
}

type Printer struct {
	doc           strings.Builder
	lineComments  bool
	blockComments bool
	newLines      bool
	indent        string
	indentLevel   int
	lastChar      rune
	ensureLine    bool
	ensureSpace   bool
	printer       func(*Printer, ast.Node)
}

func (p *Printer) LastChar() rune {
	return p.lastChar
}

// Init initializes the printer.
//
// Call Init before printing with Print/LnPrint/SpPrint.
// Printer middleware can be registered via UsePrinter BEFORE Init.
func (p *Printer) Init(opts ...Option) {
	cfg := &config{
		lineComments:  true,
		blockComments: true,
		newLines:      true,
		indent:        "  ",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	p.doc.Reset()
	p.lineComments = cfg.lineComments
	p.blockComments = cfg.blockComments
	p.newLines = cfg.newLines
	p.indent = cfg.indent
	p.indentLevel = 0
	p.lastChar = eol
	p.ensureLine = false
	p.ensureSpace = false
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
		case token.Token:
			p.printToken(v)
		case ast.Node:
			p.printNode(v)
		default:
			panic("Unsupported type")
		}
	}
}

func (p *Printer) LnPrint(arg any) {
	p.EnsureLine()
	p.Print(arg)
}

func (p *Printer) SpPrint(arg any) {
	p.EnsureSpace()
	p.Print(arg)
}

func (p *Printer) PrintTrivia(trivia []token.Token) {
	for _, tok := range trivia {
		switch tok.Type {
		case token.NEWLINE:
			if !p.newLines {
				continue
			}
			p.writeRune('\n')
		case token.LINE_COMMENT:
			if !p.lineComments {
				continue
			}
			p.printSpaceIfNeeded()
			p.printIndentIfNeeded()
			p.writeString("//" + tok.Literal)
		case token.BLOCK_COMMENT:
			if !p.blockComments {
				continue
			}
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
	if p.newLines && !isNewLine(p.lastChar) {
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
	if p.ensureSpace {
		p.printSpaceIfNeeded()
		p.ensureSpace = false
		p.ensureLine = false
	}
	if p.ensureLine {
		p.printLineIfNeeded()
		p.ensureLine = false
	}
	p.printIndentIfNeeded()
}
