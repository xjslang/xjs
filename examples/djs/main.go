package main

import (
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var deferTyp = token.RegisterType("defer")

type DeferStmt struct {
	ast.BaseStmt
	DeferToken token.Token
	Stmt       ast.Stmt
}

func djsPlugin(b *builder.Builder) {
	// the scanner that can read "defer"
	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "defer" {
			tok.Type = deferTyp
		}
		return tok
	})
	// the parser can now parse "defer"
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (node ast.Stmt, err error) {
		if p.CurrentToken.Type == deferTyp {
			deferStmt := &DeferStmt{DeferToken: p.CurrentToken}
			p.AdvanceToken() // consume "defer"
			if deferStmt.Stmt, err = js.ParseStmt(p); err != nil {
				return
			}
			node = deferStmt
			return
		}
		return next()
	})
}

func main() {
	input := `
	function foo() {
		// ensures closing db properly
		let db = openDb()
		defer closeDb()

		// ensures closing file properly
		let file = openFile()
		defer {
			print('Closing file...')
			closeFile()
		}
	}`

	djsParser := xjs.NewBuilder().
		Install(djsPlugin).
		Build([]byte(input))
	node, err := js.ParseProgram(djsParser)
	if err != nil {
		panic(err)
	}

	// create a compiler that can compile `DeferStmt`
	djsCompiler := xjs.NewPrinter()
	djsCompiler.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(ast.Node)) {
		if node, ok := node.(*DeferStmt); ok {
			pr.PrintTrivia(node.DeferToken.LeadingTrivia) // print previous comments and new lines
			pr.LnPrint("{using _ = {[Symbol.dispose]() {")
			pr.Print(node.Stmt)
			pr.Print("}}}")
			return
		}
		next(node)
	})
	djsCompiler.Init()
	djsCompiler.Print(node)
	fmt.Println(djsCompiler.String())

	// create a formatter that can format `DeferStmt`
	djsFormatter := xjs.NewPrinter()
	djsFormatter.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(ast.Node)) {
		if node, ok := node.(*DeferStmt); ok {
			pr.EnsureLine() // ensure a new line is added before printing
			pr.Print(node.DeferToken)
			pr.EnsureSpace() // ensure a new space is added before printing
			if deferNode, ok := node.Stmt.(*js.ExprStmt); ok {
				pr.Print(deferNode.Expr)
				pr.Print(deferNode.Layout.Semi)
			} else {
				pr.Print(node.Stmt)
			}
			return
		}
		next(node)
	})
	djsFormatter.Init()
	djsFormatter.Print(node)
	fmt.Println(djsFormatter.String())
}
