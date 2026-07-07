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

// TODO: the compact option should ignore spaces introduced by the `Space()` function
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
	ensureChar   rune
	ensure       bool
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
	p.ensureChar = eol
	p.ensure = false
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

// Ensure ensures that a character is printed before printing the next text.
//
// Only one ensure character can be pending at a time. For example,
// Ensure(' ').Ensure('\n') will print a space and discard the newline request.
func (p *Printer) Ensure(c rune) *Printer {
	if p.ensure {
		return p
	}
	p.ensureChar = c
	p.ensure = true
	return p
}

// Beside ensures that the next text to print appears "beside" the previous text.
func (p *Printer) Beside() *Printer {
	return p.Ensure(eol)
}

// Line ensures that a newline is printed before printing the next text.
//
// It does not print a newline immediately, but "ensures" that a newline is printed
// between the current print and the next print.
func (p *Printer) Line() *Printer {
	return p.Ensure('\n')
}

// Space ensures that a space is printed before printing the next text.
//
// It does not print a space immediately, but "ensures" that a space is printed
// between the current print and the next print "on the same line".
func (p *Printer) Space() *Printer {
	return p.Ensure(' ')
}

func (p *Printer) Print(args ...any) *Printer {
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
	return p
}

func (p *Printer) PrintTrivia(trivia []token.Token) {
	es, e := p.ensureChar, p.ensure
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
	p.ensureChar, p.ensure = es, e
}

func (p *Printer) Error(msg string) error {
	s := p.doc.String()
	line, col := 0, 0
	for _, c := range s {
		if c == '\n' {
			line++
			col = -1
		}
		col++
	}
	return ErrorAt(token.Position{
		Line:   line,
		Column: col,
	}, msg)
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
	if p.ensure {
		switch p.ensureChar {
		case '\n':
			if p.withNewLines && !isNewLine(p.lastChar) {
				p.writeRune(p.ensureChar)
			}
		case ' ':
			if !isWhitespace(p.lastChar) {
				p.writeRune(' ')
			}
		case eol:
		default:
			if p.lastChar != p.ensureChar {
				p.writeRune(p.ensureChar)
			}
		}
		p.ensureChar = eol
		p.ensure = false
	}
	p.printIndentIfNeeded()
}
