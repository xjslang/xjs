package ast

import (
	"strings"

	"github.com/xjslang/xjs/sourcemap"
)

type CodeWriter struct {
	strings.Builder
	Mapper *sourcemap.SourceMapper
}

func (w *CodeWriter) WriteString(s string) {
	w.Builder.WriteString(s)
	if w.Mapper == nil {
		return
	}
	w.Mapper.AdvanceString(s)
}

func (w *CodeWriter) WriteRune(r rune) {
	w.Builder.WriteRune(r)
	if w.Mapper == nil {
		return
	}
	if r == '\n' {
		w.Mapper.AdvanceLine()
	} else {
		w.Mapper.AdvanceColumn(1)
	}
}

func (w *CodeWriter) AddMapping(sourceLine, sourceColumn int) {
	if w.Mapper == nil {
		return
	}
	w.Mapper.AddMapping(sourceLine, sourceColumn)
}

func (w *CodeWriter) AddNamedMapping(sourceLine, sourceColumn int, name string) {
	if w.Mapper == nil {
		return
	}
	w.Mapper.AddNamedMapping(sourceLine, sourceColumn, name)
}
