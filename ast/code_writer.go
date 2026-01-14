package ast

import (
	"strings"

	"github.com/xjslang/xjs/sourcemap"
	"github.com/xjslang/xjs/token"
)

type CodeWriter struct {
	Builder         strings.Builder
	Mapper          *sourcemap.SourceMapper
	PrettyPrint     bool
	IndentLevel     int
	IndentString    string
	WriteSemicolons bool
}

func (cw *CodeWriter) String() string {
	return cw.Builder.String()
}

func (cw *CodeWriter) WriteString(s string) {
	cw.Builder.WriteString(s)
	if cw.Mapper == nil {
		return
	}
	cw.Mapper.AdvanceString(s)
}

func (cw *CodeWriter) WriteRune(r rune) {
	cw.Builder.WriteRune(r)
	if cw.Mapper == nil {
		return
	}
	if r == '\n' {
		cw.Mapper.AdvanceLine()
	} else {
		cw.Mapper.AdvanceColumn(1)
	}
}

func (cw *CodeWriter) AddMapping(pos token.Position) {
	if cw.Mapper == nil {
		return
	}
	cw.Mapper.AddMapping(pos.Line, pos.Column)
}

func (cw *CodeWriter) AddNamedMapping(sourceLine, sourceColumn int, name string) {
	if cw.Mapper == nil {
		return
	}
	cw.Mapper.AddNamedMapping(sourceLine, sourceColumn, name)
}

// WriteSemi writes a semicolon if WriteSemicolons is true.
func (cw *CodeWriter) WriteSemi() {
	if !cw.PrettyPrint {
		cw.WriteRune(';')
		return
	}

	// In pretty-print mode, respect WriteSemicolons setting
	if cw.WriteSemicolons {
		cw.WriteRune(';')
	}
}

// WriteIndent writes the current indentation level
func (cw *CodeWriter) WriteIndent() {
	if !cw.PrettyPrint {
		return
	}
	indent := cw.IndentString
	if indent == "" {
		indent = "  " // default: 2 spaces
	}
	for i := 0; i < cw.IndentLevel; i++ {
		cw.WriteString(indent)
	}
}

// WriteNewline writes a newline character if PrettyPrint is enabled
func (cw *CodeWriter) WriteNewline() {
	if !cw.PrettyPrint {
		return
	}
	cw.WriteRune('\n')
}

// WriteSpace writes a space character if PrettyPrint is enabled
func (cw *CodeWriter) WriteSpace() {
	if !cw.PrettyPrint {
		return
	}
	cw.WriteRune(' ')
}

// IncreaseIndent increases the indentation level
func (cw *CodeWriter) IncreaseIndent() {
	if !cw.PrettyPrint {
		return
	}
	cw.IndentLevel++
}

// DecreaseIndent decreases the indentation level
func (cw *CodeWriter) DecreaseIndent() {
	if !cw.PrettyPrint {
		return
	}
	if cw.IndentLevel > 0 {
		cw.IndentLevel--
	}
}
