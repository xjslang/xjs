package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
)

func Parse(input []byte) (*js.Program, error) {
	sc := ScannerBuilder().Build(input)
	p := ParserBuilder().Build(sc)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := PrinterBuilder().Build(opts...)
	pr.Print(result)
	return pr.Output()
}

func ScannerBuilder() *scanner.Builder {
	return js.ScannerBuilder()
}

func ParserBuilder() *parser.Builder {
	return js.ParserBuilder()
}

func PrinterBuilder() *printer.Builder {
	return js.PrinterBuilder()
}
