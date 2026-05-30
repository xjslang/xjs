package xjs_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

func Example_basic() {
	input := `function hello() {
	let x = 100
	let y = 200
}`
	s := xjs.NewScanner()
	s.Init([]byte(input))
	p := xjs.NewParser()
	p.Init(s)
	result, err := p.Parse()
	if err != nil {
		panic(err)
	}

	pr := xjs.NewPrinter()
	pr.Init()
	pr.Print(result)
	fmt.Print(pr.String())
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}

type iifeExpr struct {
	LparenToken token.Token
	RparenToken token.Token
	Function    *js.FunctionDecl
}

func (node *iifeExpr) Type() string {
	return "iifeExpr"
}

func TestMiddlewares(t *testing.T) {
	input := `(function foo() {
		print('Hello, World!')
	})()`
	s := xjs.NewScanner()
	s.Init([]byte(input))
	p := xjs.NewParser()
	// parse IIFE expressions
	p.UsePrefixParser(func(p *parser.Parser, next func() (ast.Node, error)) (_ ast.Node, err error) {
		if p.CurrentToken.Type == token.LPAREN && p.PeekToken.Type == js.FUNCTION {
			node := &iifeExpr{LparenToken: p.CurrentToken}
			p.AdvanceToken()
			if node.Function, err = js.ParseFunctionDecl(p); err != nil {
				return
			}
			if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
				return
			}
			return node, nil
		}
		return next()
	})
	p.Init(s)
	result, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	pr := xjs.NewPrinter()
	pr.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		if node, ok := node.(*iifeExpr); ok {
			p.Print(node.LparenToken)
			p.Print(node.Function)
			p.Print(node.RparenToken)
			return
		}
		next(node)
	})
	pr.Init(printer.WithIndent("\t"))
	pr.Print(result)
	expected := "(\nfunction foo() {\n\tprint('Hello, World!');\n})();"
	require.Equal(t, expected, pr.String())
}
