package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type Scanner struct {
	scanner.Scanner
}

func (s *Scanner) UseCoreScanners() {
	s.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "function" {
			tok.Type = js.FUNCTION
		}
		return tok
	})
}

// TODO: https://github.com/xjslang/xjs/pull/141#discussion_r3321072513
func (s *Scanner) Init(input []byte) {
	s.UseCoreScanners()
	s.Scanner.Init(input)
}

type Parser struct {
	parser.Parser
}

func (p *Parser) UseCoreParsers() {
	p.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == js.FUNCTION {
			return js.ParseFunction(p)
		}
		return next()
	})
}

// TODO: https://github.com/xjslang/xjs/pull/141#discussion_r3321072535
func (p *Parser) Init(s parser.Scanner) {
	p.UseCoreParsers()
	p.Parser.Init(s)
}

type Printer struct {
	printer.Printer
}

func (p *Printer) UseCorePrinters() {
	p.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		if node, ok := node.(*js.Function); ok {
			js.PrintFunction(p, node)
			return
		}
		next(node)
	})
}

// TODO: https://github.com/xjslang/xjs/pull/141#discussion_r3321072548
func (p *Printer) Init(opts ...printer.Option) {
	p.UseCorePrinters()
	p.Printer.Init(opts...)
}
