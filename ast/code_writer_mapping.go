package ast

import "github.com/xjslang/xjs/token"

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
