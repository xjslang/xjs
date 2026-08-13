package xjs_test

import (
	"fmt"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/jsextended"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var PIPE = token.RegisterType("PIPE", "|>")

type InfixPipeOp struct {
	ast.BaseExpr // operands are expressions
	Left, Right  ast.Expr
}

// Example_infixOp_pipeline demonstrates how to implement a custom infix
// operator that introduces a pipeline syntax (|>) to JavaScript expressions.
func Example_infixOp_pipeline() {
	// Don't forget to register PIPE as an infix operator!
	token.RegisterInfixOp(PIPE, token.COMMA.Precedence())

	// Create a scanner that recognize the "|>" token.
	sb := jsextended.ScannerBuilder()
	sb.UseScanner(func(sc *scanner.Scanner, next func(*scanner.Scanner) (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(sc); err != nil {
			return
		}
		if tok.Literal == "|" && sc.CurrentChar() == '>' {
			sc.AdvanceChar()
			tok.Type = PIPE
			tok.Literal = "|>"
		}
		return
	})

	// Create a parser that "understand" the pipeline syntax.
	pb := jsextended.ParserBuilder()
	pb.UseInfixOpParser(func(p *parser.Parser, left ast.Expr, next func(*parser.Parser, ast.Expr) (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type != PIPE {
			return next(p, left) // delegate to the "next" middleware
		}
		p.AdvanceToken()
		node := &InfixPipeOp{Left: left}
		// Unlike ParseExpr, ParseRightExpr stops when it finds an operator of equal or lower precedence.
		// That is, if the parser finds another "|>", it will stop parsing.
		if node.Right, err = js.ParseRightExpr(p, PIPE.Precedence()); err != nil {
			return
		}
		return node, nil
	})

	// Create a printer able to print the pipeline node.
	prb := jsextended.PrinterBuilder()
	prb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
		switch v := node.(type) {
		case *InfixPipeOp:
			ctx := pr.PushContext()
			defer pr.PopContext()
			ctx["it"] = v.Left
			if w, ok := v.Left.(*js.Variable); ok {
				// TODO: relocate comments to the right member
				w.LeadingTrivia = ""
			}
			pr.Print(v.Right)
			return nil
		case *js.Variable:
			if v.Literal == "it" {
				ctx := pr.Context()
				if _, ok := ctx["it"]; ok {
					pr.Print("(", ctx["it"], ")")
					return nil
				}
			}
			pr.Print(v.Token)
			return nil
		}
		return next(pr, node)
	})

	// parsing
	input := `const result = raw |> it.trim() |> it.toLowerCase() |> encodeURIComponent(it);`
	s := sb.Build(input)
	p := pb.Build(s)
	result, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}

	// printing
	pr := prb.Build()
	pr.Print(result)
	code, err := pr.Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(code)
	// Output:
	// const result = encodeURIComponent((((raw).trim()).toLowerCase()));
}
