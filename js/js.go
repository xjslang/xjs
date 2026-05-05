package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
)

func Parse(input []byte) (ast.Statement, error) {
	l := &lexer.Lexer{}
	l.Init(input)
	b := Builder{}
	b.InstallCorePlugins()
	p := b.Build(l)
	return ParseProgram(p)
}
