package ast

import (
	"strings"

	"github.com/xjslang/xjs/sourcemap"
)

type CodeWriter struct {
	Builder         strings.Builder
	Mapper          *sourcemap.SourceMapper
	PrettyPrint     bool
	IndentLevel     int
	IndentString    string
	WriteSemicolons bool

	pendings []rune
}

// WriteString writes a string to the buffer
func (cw *CodeWriter) WriteString(s string) {
	cw.flushPending()
	cw.Builder.WriteString(s)
	if cw.Mapper == nil {
		return
	}
	cw.Mapper.AdvanceString(s)
}

// WriteRune writes a rune to the buffer
func (cw *CodeWriter) WriteRune(r rune) {
	cw.flushPending()
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

// String returns the accumulated string
func (cw *CodeWriter) String() string {
	return cw.Builder.String()
}
