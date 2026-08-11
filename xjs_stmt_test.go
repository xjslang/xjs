package xjs_test

import (
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var DEFER = token.RegisterType("DEFER", "defer")

type DeferStmt struct {
	ast.BaseStmt
	Defer token.Token // preserve defer token with its leading trivia
	Stmt  ast.Stmt
}

// Example_stmt_defer demonstrates how to extend xjs with a custom statement
// by adding a Go-style defer statement to JavaScript.
func Example_stmt_defer() {
	// Create a scanner that recognizes the "defer" token.
	sb := xjs.ScannerBuilder()
	sb.UseScanner(func(sc *scanner.Scanner, next func(*scanner.Scanner) (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(sc); err != nil {
			return
		}
		if tok.Literal == "defer" {
			tok.Type = DEFER
		}
		return
	})

	// Create a parser that "understand" the defer syntax.
	pb := xjs.ParserBuilder()
	pb.UseStmtParser(func(p *parser.Parser, next func(*parser.Parser) (ast.Stmt, error)) (node ast.Stmt, err error) {
		if p.CurrentToken.Type != DEFER {
			return next(p) // delegate to the "next" middleware
		}
		deferStmt := &DeferStmt{Defer: p.CurrentToken}
		p.AdvanceToken() // consume "defer"
		if deferStmt.Stmt, err = js.ParseStmt(p); err != nil {
			return
		}
		node = deferStmt
		return
	})

	// Create a compiler able to compile the defer node.
	cb := xjs.PrinterBuilder()
	cb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
		switch v := node.(type) {
		case *DeferStmt:
			pr.PrintTrivia(v.Defer.LeadingTrivia) // print previous comments and new lines
			pr.Line().Print("{using _ = {[Symbol.dispose]() {")
			pr.Print(v.Stmt)
			pr.Print("}}}")
			return nil
		}
		return next(pr, node)
	})

	// create a formatter that can format `DeferStmt`
	fb := xjs.PrinterBuilder()
	fb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
		if node, ok := node.(*DeferStmt); ok {
			pr.Line().Print(node.Defer)
			pr.Space() // ensure a new space is added
			if deferNode, ok := node.Stmt.(*js.ExprStmt); ok {
				pr.Print(deferNode.Expr)
			} else {
				pr.Print(node.Stmt)
			}
			return nil
		}
		return next(pr, node)
	})

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

	// parsing
	s := sb.Build([]byte(input))
	p := pb.Build(s)
	node, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}

	// compiling
	c := cb.Build()
	c.Print(node)
	compiled, err := c.Output()
	if err != nil {
		panic(err)
	}

	// formatting
	f := fb.Build()
	f.Print(node)
	formatted, err := f.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(compiled)
	fmt.Println(formatted)
	// Output:
	// function foo() {
	//   // ensures closing db properly
	//   let db = openDb();
	//   {using _ = {[Symbol.dispose]() {
	//   closeDb();}}}
	//
	//   // ensures closing file properly
	//   let file = openFile();
	//   {using _ = {[Symbol.dispose]() {{
	//     print('Closing file...');
	//     closeFile();
	//   }}}}
	// }
	//
	// function foo() {
	//   // ensures closing db properly
	//   let db = openDb();
	//   defer closeDb()
	//
	//   // ensures closing file properly
	//   let file = openFile();
	//   defer {
	//     print('Closing file...');
	//     closeFile();
	//   }
	// }
}
