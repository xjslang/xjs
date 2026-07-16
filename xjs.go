package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/printer"
)

func Parse(input []byte) (*js.Program, error) {
	p := PluginBuilder().Build(input)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := PrinterBuilder().Build(opts...)
	pr.Print(result)
	return pr.Output()
}

func PluginBuilder() *plugin.Builder {
	return plugin.New().
		Install(js.Plugin).
		Install(js.ExtendedPlugin)
}

func PrinterBuilder() *printer.Builder {
	return printer.NewBuilder().
		UsePrinter(js.Printer).
		UsePrinter(js.ExtendedPrinter)
}
