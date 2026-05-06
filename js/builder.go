package js

import (
	"github.com/xjslang/xjs/parser"
)

type Builder struct {
	parser.Builder
}

func (b *Builder) InstallCorePlugins() {
	b.Install(MathPlugin)
	b.Install(LetPlugin)
	b.Install(FunctionPlugin)
}
