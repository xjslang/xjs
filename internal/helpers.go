package internal

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

func Parse(input []byte) (*js.Program, error) {
	sc := xjs.ScannerBuilder().Build(input)
	p := xjs.ParserBuilder().Build(sc)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := xjs.PrinterBuilder().Build(opts...)
	pr.Print(result)
	return pr.Output()
}
