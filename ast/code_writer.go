package ast

import (
	"strings"

	"github.com/xjslang/xjs/sourcemap"
	"github.com/xjslang/xjs/token"
)

type CodeWriter struct {
	strings.Builder
	Mapper *sourcemap.SourceMapper
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
