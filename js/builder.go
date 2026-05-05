package js

import (
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

type Builder struct {
	parser.Builder
}

func (b *Builder) Build(l *lexer.Lexer) *parser.Parser {
	b.Install(LetPlugin)
	return b.Builder.Build(l)
}
