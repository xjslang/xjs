package jsx

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

func Parse(input []byte) (*js.Program, error) {
	b := xjs.NewBuilder()
	b.Install(Plugin)
	p := b.Build(input)
	return js.ParseProgram(p)
}

func Compile(result ast.Node) (string, error) {
	pr := xjs.NewPrinter(printer.Compact())
	pr.UsePrinter(Compiler)
	pr.Print(result)
	return pr.Output()
}

func Format(result ast.Node, opts ...printer.Option) (string, error) {
	pr := xjs.NewPrinter(opts...)
	pr.UsePrinter(Formatter)
	pr.Print(result)
	return pr.Output()
}
