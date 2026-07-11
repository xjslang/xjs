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
	withLogs     bool
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

func WithLogs(value bool) Option {
	return func(cfg *config) {
		cfg.withLogs = value
	}
}

type Printer struct {
	doc          strings.Builder
	withComments bool
	withNewLines bool
	withLogs     bool
	indent       string
	indentLevel  int
	lastChar     rune
	ensureChar   rune
	ensure       bool
	printer      func(*Printer, ast.Node) error
	context      []map[string]string
	errors       ErrorList
}

func (pr *Printer) init(opts ...Option) {
	cfg := &config{
		withComments: true,
		withNewLines: true,
		indent:       "  ",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	pr.doc.Reset()
	pr.withComments = cfg.withComments
	pr.withNewLines = cfg.withNewLines
	pr.withLogs = cfg.withLogs
	pr.indent = cfg.indent
	pr.indentLevel = 0
	pr.lastChar = eol
	pr.ensureChar = eol
	pr.ensure = false
	if pr.printer == nil {
		pr.printer = defaultPrinter
	}
	pr.errors = nil
}

func (pr *Printer) IncreaseIndent() {
	pr.indentLevel++
}

func (pr *Printer) DecreaseIndent() {
	if pr.indentLevel > 0 {
		pr.indentLevel--
	}
}

func (pr *Printer) PrintIndent() {
	for range pr.indentLevel {
		pr.writeString(pr.indent)
	}
}

// Ensure ensures that a character is printed before printing the next text.
//
// Only one ensure character can be pending at a time. For example,
// Ensure(' ').Ensure('\n') will print a space and discard the newline request.
func (pr *Printer) Ensure(c rune) *Printer {
	if pr.ensure {
		return pr
	}
	pr.ensureChar = c
	pr.ensure = true
	return pr
}

// Beside ensures that the next text to print appears "beside" the previous text.
func (pr *Printer) Beside() *Printer {
	return pr.Ensure(eol)
}

// Line ensures that a newline is printed before printing the next text.
//
// It does not print a newline immediately, but "ensures" that a newline is printed
// between the current print and the next print.
func (pr *Printer) Line() *Printer {
	return pr.Ensure('\n')
}

// Space ensures that a space is printed before printing the next text.
//
// It does not print a space immediately, but "ensures" that a space is printed
// between the current print and the next print "on the same line".
func (pr *Printer) Space() *Printer {
	return pr.Ensure(' ')
}

func (pr *Printer) Print(args ...any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			pr.printString(v)
		case int:
			pr.printString(strconv.Itoa(v))
		case rune:
			pr.printRune(v)
		case token.Token:
			pr.printToken(v)
		case ast.Node:
			pr.printNode(v)
		default:
			panic("Unsupported type")
		}
	}
}

func (pr *Printer) Log(args ...any) {
	if pr.withLogs {
		pr.Print(args...)
	}
}

func (pr *Printer) PrintTrivia(trivia []token.Token) {
	es, e := pr.ensureChar, pr.ensure
	for _, tok := range trivia {
		if tok.Type == token.NEWLINE {
			if pr.withNewLines {
				pr.writeRune('\n')
			}
			continue
		}
		if pr.withComments {
			pr.printSpaceIfNeeded()
			pr.printIndentIfNeeded()
			pr.writeString(tok.Literal)
		}
	}
	pr.ensureChar, pr.ensure = es, e
}

func (pr *Printer) Error(msg string) error {
	s := pr.doc.String()
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

func (pr *Printer) Errors() ErrorList {
	return append(ErrorList{}, pr.errors...)
}

func (pr *Printer) Output() (string, error) {
	return pr.doc.String(), errors.Join(pr.errors...)
}

func (pr *Printer) writeString(s string) {
	if len(s) == 0 {
		return
	}
	r, _ := utf8.DecodeLastRuneInString(s)
	pr.lastChar = r
	pr.doc.WriteString(s)
}

func (pr *Printer) writeRune(r rune) {
	pr.lastChar = r
	pr.doc.WriteRune(r)
}

func (pr *Printer) printNode(node ast.Node) {
	if err := pr.printer(pr, node); err != nil {
		pr.errors = append(pr.errors, err)
	}
}

func (pr *Printer) printString(s string) {
	if len(s) == 0 {
		return
	}
	pr.printSeparatorIfNeeded()
	pr.writeString(s)
}

func (pr *Printer) printRune(r rune) {
	pr.printSeparatorIfNeeded()
	pr.writeRune(r)
}

func (pr *Printer) printToken(tok token.Token) {
	pr.PrintTrivia(tok.LeadingTrivia)
	pr.printString(tok.Literal)
}

func (pr *Printer) printSpaceIfNeeded() {
	if !isWhitespace(pr.lastChar) {
		pr.writeRune(' ')
	}
}

func (pr *Printer) printIndentIfNeeded() {
	if isNewLine(pr.lastChar) {
		pr.PrintIndent()
	}
}

func (pr *Printer) printSeparatorIfNeeded() {
	if pr.ensure {
		switch pr.ensureChar {
		case '\n':
			if pr.withNewLines && !isNewLine(pr.lastChar) {
				pr.writeRune(pr.ensureChar)
			}
		case ' ':
			if !isWhitespace(pr.lastChar) {
				pr.writeRune(' ')
			}
		case eol:
		default:
			if pr.lastChar != pr.ensureChar {
				pr.writeRune(pr.ensureChar)
			}
		}
		pr.ensureChar = eol
		pr.ensure = false
	}
	pr.printIndentIfNeeded()
}
