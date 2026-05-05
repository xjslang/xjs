package js

import (
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

type Builder struct {
	parser.Builder
}

// TODO: https://github.com/xjslang/xjs/pull/90#discussion_r3191186051
func (b *Builder) Build(l *lexer.Lexer) *parser.Parser {
	b.Install(LetPlugin)
	return b.Builder.Build(l)
}
