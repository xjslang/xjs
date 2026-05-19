package main

import (
	"fmt"

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
	djsScanner := &scanner.Scanner{}
	djsScanner.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "defer" {
			tok.Type = deferTyp
		}
		return tok
	})
	djsScanner.Init([]byte(input))

	// create a parser that can parse "defer"
	djsParser := &parser.Parser{}
	djsParser.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == deferTyp {
			deferStmt := &DeferStmt{DeferToken: p.CurrentToken}
			p.AdvanceToken() // consume "defer"
			deferStmt.Stmt, err = parser.ParseExprStmt(p)
			if err != nil {
				return
			}
			node = deferStmt
			return
		}
		return next()
	})
	djsParser.Init(djsScanner)
	node, err := parser.ParseProgram(djsParser)
	if err != nil {
		panic(err)
	}

	// create a compiler that can compile `DeferStmt`
	compiler := &printer.Printer{}
	compiler.UsePrinter(func(pr *printer.Printer, node ast.Node, next func()) {
		if node, ok := node.(*DeferStmt); ok {
			pr.PrintTrivia(node.DeferToken.LeadingTrivia) // print previous comments and new lines
			pr.EnsureLine()                               // ensure a new line is added before printing
			pr.PrintIndentedString("{using _ = {[Symbol.dispose]() {")
			pr.PrintNode(node.Stmt.Expr)
			pr.PrintToken(node.Stmt.SemiToken)
			pr.PrintIndentedString("}}}")
			return
		}
		next()
	})
	compiler.Init()
	compiler.PrintNode(node)
	fmt.Println(compiler.String())

	// create a formatter that can format `DeferStmt`
	fr := &printer.Printer{}
	fr.UsePrinter(func(pr *printer.Printer, node ast.Node, next func()) {
		if node, ok := node.(*DeferStmt); ok {
			pr.EnsureLine() // ensure a new line is added before printing
			pr.PrintToken(node.DeferToken)
			pr.EnsureSpace() // ensure a new space is added before printing
			pr.PrintNode(node.Stmt.Expr)
			pr.PrintToken(node.Stmt.SemiToken)
			return
		}
		next()
	})
	fr.Init()
	fr.PrintNode(node)
	fmt.Println(fr.String())
}
