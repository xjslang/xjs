package printer

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

const eol = rune(-1)

type Error struct {
	token.Position
	Message string
}

func (err Error) Error() string {
	return "[line:" + strconv.Itoa(err.Line) +
		", col:" + strconv.Itoa(err.Column) +
		"] " + err.Message
}

type ErrorList []error

func (list ErrorList) Error() string {
	return errors.Join(list...).Error()
}

type config struct {
	indent       string
	withComments bool
	withNewLines bool
}

type Option func(*config)

func Compact() Option {
	return func(cfg *config) {
		cfg.withComments = false
		cfg.withNewLines = false
	}
}

func WithIndent(value string) Option {
	return func(cfg *config) {
		cfg.indent = value
	}
}

func WithComments(value bool) Option {
	return func(cfg *config) {
		cfg.withComments = value
	}
}

func WithNewLines(value bool) Option {
	return func(cfg *config) {
		cfg.withNewLines = value
	}
}

type Printer struct {
	doc          strings.Builder
	withComments bool
	withNewLines bool
	indent       string
	indentLevel  int
	lastChar     rune
	ensureBeside bool
	ensureLine   bool
	ensureSpace  bool
	printer      func(*Printer, ast.Node) error
	context      []map[string]string
	errors       ErrorList
}

func (p *Printer) init(opts ...Option) {
	cfg := &config{
		withComments: true,
		withNewLines: true,
		indent:       "  ",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	p.doc.Reset()
	p.withComments = cfg.withComments
	p.withNewLines = cfg.withNewLines
	p.indent = cfg.indent
	p.indentLevel = 0
	p.lastChar = eol
	p.ensureBeside = false
	p.ensureLine = false
	p.ensureSpace = false
	if p.printer == nil {
		p.printer = defaultPrinter
	}
	p.errors = nil
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

// EnsureBeside ensures that the next text to print appears "beside" the previous text.
func (p *Printer) EnsureBeside() {
	p.ensureBeside = true
}

// EnsureLine ensures that a newline is printed before printing the next text.
//
// It does not print a newline inmediately, but "ensures" that a newline is printed
// between the current print and the next print.
func (p *Printer) EnsureLine() {
	p.ensureLine = true
}

// EnsureSpace ensures that a space is printed before printing the next text.
//
// It does not print a space inmediately, but "ensures" that a space is printed
// between the current print and the next print "on the same line".
func (p *Printer) EnsureSpace() {
	p.ensureSpace = true
}

func (p *Printer) Print(args ...any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			p.printString(v)
		case int:
			p.printString(strconv.Itoa(v))
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

// BsPrint ensures that the next text to print appears "beside" the previous text.
// This is a combination of EnsureBeside() + Print(a).
// It takes priority over LnPrint() and SpPrint().
func (p *Printer) BsPrint(arg any) *Printer {
	p.EnsureBeside()
	p.Print(arg)
	return p
}

// LnPrint ensures that a newline is printed before printing the next text.
// This is a combination of EnsureLine() + Print(a).
func (p *Printer) LnPrint(arg any) *Printer {
	p.EnsureLine()
	p.Print(arg)
	return p
}

// SpPrint ensures that a space is printed before printing the next text.
// This is a combination of EnsureSpace() + Print(a).
// It takes priority over LnPrint().
func (p *Printer) SpPrint(arg any) *Printer {
	p.EnsureSpace()
	p.Print(arg)
	return p
}

func (p *Printer) PrintTrivia(trivia []token.Token) {
	eb, es, el := p.ensureBeside, p.ensureSpace, p.ensureLine
	for _, tok := range trivia {
		if tok.Type == token.NEWLINE {
			if p.withNewLines {
				p.writeRune('\n')
			}
			continue
		}
		if p.withComments {
			p.printSpaceIfNeeded()
			p.printIndentIfNeeded()
			p.writeString(tok.Literal)
		}
	}
	p.ensureBeside, p.ensureSpace, p.ensureLine = eb, es, el
}

func (p *Printer) Errors() ErrorList {
	return append(ErrorList{}, p.errors...)
}

func (p *Printer) Output() (string, error) {
	return p.doc.String(), errors.Join(p.errors...)
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
	if err := p.printer(p, node); err != nil {
		p.errors = append(p.errors, err)
	}
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
	if p.withNewLines && !isNewLine(p.lastChar) {
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
	if p.ensureBeside {
		p.ensureBeside = false
		p.ensureSpace = false
		p.ensureLine = false
	}
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
