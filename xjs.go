package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

// Core scanners.
var CoreScanners = []func(*scanner.Scanner, func() token.Token) token.Token{
	func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "function" {
			tok.Type = js.FUNCTION
		}
		return tok
	},
}

// Core statement parsers.
var CoreStmtParsers = []func(*parser.Parser, func() (ast.Node, error)) (ast.Node, error){
	func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == js.FUNCTION {
			return js.ParseFunction(p)
		}
		return next()
	},
}

// Core printers.
var CorePrinters = []func(*printer.Printer, ast.Node, func(node ast.Node)){
	func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		if node, ok := node.(*js.Function); ok {
			js.PrintFunction(p, node)
			return
		}
		next(node)
	},
}

func NewScanner() *scanner.Scanner {
	s := &scanner.Scanner{}
	// use middlewares
	for _, md := range CoreScanners {
		s.UseScanner(md)
	}
	return s
}

func NewParser() *parser.Parser {
	p := &parser.Parser{}
	// use middlewares
	for _, md := range CoreStmtParsers {
		p.UseStmtParser(md)
	}
	return p
}

func NewPrinter() *printer.Printer {
	p := &printer.Printer{}
	// use middlewares
	for _, md := range CorePrinters {
		p.UsePrinter(md)
	}
	return p
}
