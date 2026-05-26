package formatter

import (
	"strings"
	"unicode/utf8"
)

const EOL = rune(-1) // end of line

type formatterConfig struct {
	indent string
}

type FormatterOption func(*formatterConfig)

func WithIndent(value string) FormatterOption {
	return func(cfg *formatterConfig) {
		cfg.indent = value
	}
}

type Formatter struct {
	doc         strings.Builder
	indent      string
	indentLevel int
	LastChar    rune
}

func (f *Formatter) Init(opts ...FormatterOption) {
	cfg := &formatterConfig{
		indent: "  ",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	f.doc.Reset()
	f.indent = cfg.indent
	f.indentLevel = 0
	f.LastChar = EOL
}

func (f *Formatter) PrintString(s string) {
	if len(s) == 0 {
		return
	}
	r, _ := utf8.DecodeLastRuneInString(s)
	f.LastChar = r
	f.doc.WriteString(s)
}

func (f *Formatter) PrintRune(r rune) {
	f.LastChar = r
	f.doc.WriteRune(r)
}

func (f *Formatter) IncreaseIndent() {
	f.indentLevel++
}

func (f *Formatter) DecreaseIndent() {
	if f.indentLevel > 0 {
		f.indentLevel--
	}
}

func (f *Formatter) PrintIndent() {
	for range f.indentLevel {
		f.PrintString(f.indent)
	}
}

func (f *Formatter) String() string {
	return f.doc.String()
}

func (f *Formatter) Bytes() []byte {
	return []byte(f.String())
}
