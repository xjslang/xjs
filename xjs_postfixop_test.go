package xjs_test

import (
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type PostfixFactorialOp struct {
	ast.BaseExpr // operands are expressions
	Value        ast.Expr
}

// Example_postfixOp_factorial demonstrates how to add a postfix "!" factorial operator.
func Example_postfixOp_factorial() {
	// Register "!" as an infix operator.
	//
	// Note: "!" is already the logical negation prefix operator. By registering
	// it as an infix operator we extend its functionality.
	token.RegisterInfixOp(token.NOT, token.LPAREN.Precedence())

	// Create a parser that recognizes the factorial syntax.
	pb := xjs.ParserBuilder()
	pb.UseInfixOpParser(func(p *parser.Parser, left ast.Expr, next func(*parser.Parser, ast.Expr) (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type != token.NOT {
			return next(p, left)
		}
		p.AdvanceToken()
		node := &PostfixFactorialOp{Value: left}
		// Note: a postfix operator can be considered as an infix operator without a right-hand operand.
		return node, nil
	})

	// Create a printer able to print the factorial node.
	cb := xjs.PrinterBuilder()
	cb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
		switch v := node.(type) {
		case *PostfixFactorialOp:
			pr.Print("(function fact(n) { return n <= 1 ? 1 : n * fact(n - 1); })(", v.Value, ")")
			return nil
		}
		return next(pr, node)
	})

	// parsing
	input := `let f = 5!`
	s := xjs.ScannerBuilder().Build([]byte(input))
	p := pb.Build(s)
	result, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}

	// compiling
	c := cb.Build()
	c.Print(result)
	compiled, err := c.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(compiled)
	// Output:
	// let f = (function fact(n) { return n <= 1 ? 1 : n * fact(n - 1); })(5);
}
