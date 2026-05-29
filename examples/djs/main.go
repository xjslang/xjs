package main

import (
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var deferTyp = token.RegisterType("defer")

// implements ast.Node
type DeferStmt struct {
	DeferToken token.Token
	Stmt       *ast.ExprStmt
}

func (node *DeferStmt) Type() string {
	return "DeferStmt"
}

func main() {
	input := `
	function foo() {
		let db = openDb()
		// ensures closing db properly
		defer closeDb()
	}`

	// create a scanner that can read "defer"
	djsScanner := xjs.NewScanner()
	djsScanner.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "defer" {
			tok.Type = deferTyp
		}
		return tok
	})
	djsScanner.Init([]byte(input))

	// create a parser that can parse "defer"
	djsParser := xjs.NewParser()
	djsParser.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == deferTyp {
			deferStmt := &DeferStmt{DeferToken: p.CurrentToken}
			p.AdvanceToken() // consume "defer"
			if deferStmt.Stmt, err = parser.ParseExprStmt(p); err != nil {
				return
			}
			node = deferStmt
			return
		}
		return next()
	})
	djsParser.Init(djsScanner)
	node, err := djsParser.Parse()
	if err != nil {
		panic(err)
	}

	// create a compiler that can compile `DeferStmt`
	djsCompiler := xjs.NewPrinter()
	djsCompiler.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(ast.Node)) {
		if node, ok := node.(*DeferStmt); ok {
			pr.PrintTrivia(node.DeferToken.LeadingTrivia) // print previous comments and new lines
			pr.LnPrint("{using _ = {[Symbol.dispose]() {")
			pr.Print(node.Stmt.Expr)
			pr.Print(node.Stmt.SemiToken)
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
			pr.Print(node.Stmt.Expr)
			pr.Print(node.Stmt.SemiToken)
			return
		}
		next(node)
	})
	djsFormatter.Init()
	djsFormatter.Print(node)
	fmt.Println(djsFormatter.String())
}
