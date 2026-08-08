// Package xjs provides tools for scanning, parsing, and printing JavaScript.
//
// It is organized in two main subpackages:
//
//   - js: a simplified subset of JavaScript, without arrow functions,
//     try/catch, or other additional ECMAScript features.
//   - jsextended: the "extended" version with additional syntax
//     and features.
//
// The exported functions in this package are thin wrappers around the
// equivalent functions in the js package.
package xjs

import (
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
)

// ScannerBuilder extends the simplified JavaScript scanner, allowing the input text to be split into custom tokens.
func ScannerBuilder() *scanner.Builder {
	return js.ScannerBuilder()
}

// ParserBuilder extends the simplified JavaScript parser, allowing custom language features to be parsed.
func ParserBuilder() *parser.Builder {
	return js.ParserBuilder()
}

// PrinterBuilder extends the simplified JavaScript printer, allowing custom language features to be printed.
func PrinterBuilder() *printer.Builder {
	return js.PrinterBuilder()
}
