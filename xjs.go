package xjs

import (
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
)

func ScannerBuilder() *scanner.Builder {
	return js.ScannerBuilder()
}

func ParserBuilder() *parser.Builder {
	return js.ParserBuilder()
}

func PrinterBuilder() *printer.Builder {
	return js.PrinterBuilder()
}
