package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

func NewBuilder() *builder.Builder {
	return builder.New().
		Install(js.FuncDeclPlugin).
		Install(js.LetStmtPlugin).
		Install(js.CallExprPlugin)
}

func NewPrinter() *printer.Printer {
	p := &printer.Printer{}
	p.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		switch v := node.(type) {
		case *js.FunctionDecl:
			js.PrintFunctionDecl(p, v)
			return
		case *js.LetStmt:
			js.PrintLetStmt(p, v)
			return
		case *js.CallExpr:
			js.PrintCallExpr(p, v)
			return
		}
		next(node)
	})
	return p
}

func Parse(input []byte) (*ast.Program, error) {
	return NewBuilder().Build(input).Parse()
}
