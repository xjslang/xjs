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

var (
	START_TAG = token.RegisterType("START_TAG", "<")
	END_TAG   = token.RegisterType("END_TAG", "</")
)

func Example_example1() {
	sb := xjs.ScannerBuilder()
	sb.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		switch tok.Type {
		case token.LT: // <
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
	})

	input := `<p>"Hello, World!"</p>`
	sc := sb.Build([]byte(input))
	for tok := sc.NextToken(); tok.Type != token.EOF; tok = sc.NextToken() {
		fmt.Printf("[%q:%s] ", tok.Literal, tok.Type.Name())
	}
	// output: ["<":START_TAG] ["p":IDENT] [">":GT] ["\"Hello, World!\"":STRING] ["</":END_TAG] ["p":IDENT] [">":GT]
}

type PrefixTagOp struct {
	ast.Expr
	Name     *js.Ident
	Children []ast.Expr
}

func Example_example2() {
	token.RegisterUnaryType(START_TAG)

	sb := xjs.ScannerBuilder()
	sb.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		switch tok.Type {
		case token.LT: // <
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
	})

	pb := xjs.ParserBuilder()
	pb.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (node ast.Expr, err error) {
		if p.CurrentToken.Type == START_TAG {
			p.AdvanceToken()
			tagNode := &PrefixTagOp{}
			if tagNode.Name, err = js.ParseIdent(p); err != nil {
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
				tagNode.Children = append(tagNode.Children, child)
			}
			if _, err = p.Expect(END_TAG); err != nil {
				return
			}
			var ident *js.Ident
			if ident, err = js.ParseIdent(p); err != nil {
				return
			}
			if tagNode.Name.Literal != ident.Literal {
				err = p.ErrorAt(ident.Token, "</"+tagNode.Name.Literal+"> expected")
				return
			}
			if _, err = p.Expect(token.GT); err != nil {
				return
			}
			return tagNode, nil
		}
		return next()
	})

	input := `let node = <p>"Hello, World!"</p>`
	sc := sb.Build([]byte(input))
	p := pb.Build(sc)
	result, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}

	code, err := xjs.Print(result)
	if err != nil {
		panic(err)
	}
	fmt.Println(code)
	// Output: let node = <unknown>;
}

func Example_example3() {
	token.RegisterUnaryType(START_TAG)

	sb := xjs.ScannerBuilder()
	sb.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		switch tok.Type {
		case token.LT: // <
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
	})

	pb := xjs.ParserBuilder()
	pb.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (node ast.Expr, err error) {
		if p.CurrentToken.Type == START_TAG {
			p.AdvanceToken()
			tagNode := &PrefixTagOp{}
			if tagNode.Name, err = js.ParseIdent(p); err != nil {
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
				tagNode.Children = append(tagNode.Children, child)
			}
			if _, err = p.Expect(END_TAG); err != nil {
				return
			}
			var ident *js.Ident
			if ident, err = js.ParseIdent(p); err != nil {
				return
			}
			if tagNode.Name.Literal != ident.Literal {
				err = p.ErrorAt(ident.Token, "</"+tagNode.Name.Literal+"> expected")
				return
			}
			if _, err = p.Expect(token.GT); err != nil {
				return
			}
			return tagNode, nil
		}
		return next()
	})

	input := `let node = <p>"Hello, World!"</p>`
	sc := sb.Build([]byte(input))
	p := pb.Build(sc)
	result, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}

	prb := xjs.PrinterBuilder()
	prb.UsePrinter(func(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		switch v := node.(type) {
		case *PrefixTagOp:
			pr.Space().Print("<", v.Name, ">")
			pr.IncreaseIndent()
			for _, child := range v.Children {
				pr.Print(child)
			}
			pr.DecreaseIndent()
			pr.Print("</", v.Name, ">")
		default:
			return next(node)
		}
		return nil
	})

	pr := prb.Build()
	pr.Print(result)
	code, err := pr.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(code)
	// Output: let node = <p>"Hello, World!"</p>;
}
