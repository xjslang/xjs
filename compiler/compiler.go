package compiler

import (
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/sourcemap"
)

type CompileResult struct {
	Code      string
	SourceMap *sourcemap.SourceMap
}

type Compiler struct {
	generateSourceMap bool
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) WithSourceMap() *Compiler {
	c.generateSourceMap = true
	return c
}

func (c *Compiler) Compile(program *ast.Program) CompileResult {
	w := ast.CodeWriter{Builder: strings.Builder{}}
	if c.generateSourceMap {
		w.Mapper = sourcemap.New()
	}
	program.WriteTo(&w)

	code := w.String()
	var sm *sourcemap.SourceMap
	if c.generateSourceMap {
		sm = w.Mapper.SourceMap()
	}

	return CompileResult{Code: code, SourceMap: sm}
}
