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

// Register our custom tokens first.
var (
	START_TAG = token.RegisterType("START_TAG", "<")
	END_TAG   = token.RegisterType("END_TAG", "</")
)

// The START_TAG token is treated as a "prefix operator".
type PrefixTagOp struct {
	ast.BaseExpr // must extend ast.BaseExpr, since operators are expressions
	Name         *js.Ident
	Children     []ast.Expr
}

// "Teach" the scanner how to recognize our custom tokens.
func htmlScanner(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
	if tok, err = next(); err != nil {
		return
	}
	switch tok.Type {
	case token.LT:
		c := sc.CurrentChar()
		switch {
		case scanner.IsLetter(c):
			tok.Type = START_TAG
		case c == '/':
			sc.AdvanceChar()
			tok.Type = END_TAG
			tok.Literal = "</"
		}
	}
	return
}

// "Teach" the parser how to traverse our custom syntax.
func htmlPrefixParser(p *parser.Parser, next func() (ast.Expr, error)) (_ ast.Expr, err error) {
	if p.CurrentToken.Type != START_TAG {
		return next() // delegate to the "next" parser
	}
	node := &PrefixTagOp{}
	if _, err = p.Expect(START_TAG); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	if _, err = p.Expect(token.GT); err != nil {
		return
	}
	for p.CurrentToken.Type != END_TAG {
		var child ast.Expr
		if child, err = p.ParseExpr(); err != nil {
			return
		}
		node.Children = append(node.Children, child)
	}
	if _, err = p.Expect(END_TAG); err != nil {
		return
	}
	var ident *js.Ident
	if ident, err = js.ParseIdent(p); err != nil {
		return
	}
	if ident.Literal != node.Name.Literal {
		return nil, p.ErrorAt(
			ident.Token,
			"Expected </"+node.Name.Literal+">",
		)
	}
	if _, err = p.Expect(token.GT); err != nil {
		return
	}
	return node, nil
}

// "Teach" the printer how to compile our custom nodes.
func htmlCompiler(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *PrefixTagOp:
		pr.Print("(function(){")
		pr.Print("const elem = document.createElement('", v.Name, "');")
		for _, child := range v.Children {
			pr.Print("elem.append(", child, ");")
		}
		pr.Print("return elem;})()")
	default:
		return next(node)
	}
	return nil
}

// Example_prefixOp demonstrates how to extend XJS with a custom prefix operator
// by adding a minimal HTML-like tag syntax (<p>...</p>) inside JavaScript expressions.
func Example_prefixOp_htmlTag() {
	source := `let x = <p>"Hello, World!"</p>`

	// Don't forget to register your prefix/infix operators!
	token.RegisterPrefixOp(START_TAG)

	// Scanning
	//
	// create a scanner that splits
	// the source into tokens
	sb := xjs.ScannerBuilder()
	sb.UseScanner(htmlScanner)
	sc := sb.Build([]byte(source)) // scanner

	// Parsing
	//
	// the parser uses the scanner above
	// to process the tokens and produce the AST
	pb := xjs.ParserBuilder()
	pb.UsePrefixOpParser(htmlPrefixParser)
	p := pb.Build(sc)                 // scanner from above
	result, err := js.ParseProgram(p) // produces the AST
	if err != nil {
		panic(err)
	}

	// Printing
	//
	// the printer turns AST nodes
	// back into text
	//
	// Different printers can be used
	// for different purposes (compile, format, etc.)
	prb := xjs.PrinterBuilder()
	prb.UsePrinter(htmlCompiler)
	compiler := prb.Build(printer.Compact())
	compiler.Print(result)         // prints the AST node
	code, err := compiler.Output() // get the output
	if err != nil {
		panic(err)
	}
	fmt.Println(code)

	// Output:
	// let x = (function(){const elem = document.createElement('p');elem.append("Hello, World!");return elem;})();
}
