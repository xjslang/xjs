package printer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type AsyncFunctionDecl struct {
	*js.FunctionDecl
	Layout struct {
		Async token.Token
	}
}

type AwaitStmt struct {
	ast.BaseStmt
	Layout struct {
		Await token.Token
	}
	Expr ast.Expr
}

func TestPrinterContext(t *testing.T) {
	asyncTyp := token.RegisterType("async")
	awaitTyp := token.RegisterType("await")

	b := xjs.PluginBuilder()
	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		if tok.Type == token.IDENT {
			switch tok.Literal {
			case "async":
				tok.Type = asyncTyp
			case "await":
				tok.Type = awaitTyp
			}
		}
		return
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (_ ast.Stmt, err error) {
		switch p.CurrentToken.Type {
		case asyncTyp:
			node := &AsyncFunctionDecl{}
			node.Layout.Async = p.CurrentToken
			p.AdvanceToken() // consume "async"
			if node.FunctionDecl, err = js.ParseFunctionDecl(p); err != nil {
				return
			}
			return node, nil
		case awaitTyp:
			node := &AwaitStmt{}
			node.Layout.Await = p.CurrentToken
			p.AdvanceToken() // consume "await"
			if node.Expr, err = p.ParseExpr(); err != nil {
				return
			}
			return node, nil
		}
		return next()
	})

	input := `function fetchUserData() {
		await http.fetch('/user/data')
	}`
	p := b.Build([]byte(input))
	result, err := js.ParseProgram(p)
	require.NoError(t, err)

	pr := xjs.PrinterBuilder().
		UsePrinter(func(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
			switch v := node.(type) {
			case *AsyncFunctionDecl:
				ctx := pr.PushContext()
				defer pr.PopContext()
				ctx["async"] = "yes"
				pr.LnPrint(v.Layout.Async)
				pr.SpPrint(v.FunctionDecl)
				return nil
			case *AwaitStmt:
				ctx := pr.Context()
				_, ok := ctx["async"]
				if !ok {
					return printer.ErrorAt(v.Layout.Await, "await is allowed only inside async functions")
				}
				pr.LnPrint(v.Layout.Await)
				pr.SpPrint(v.Expr)
				return nil
			}
			return next(node)
		}).
		Build()
	pr.Print(result)
	_, err = pr.Output()
	require.ErrorContains(t, err, "await is allowed only inside async functions")
	// validate pr.Errors()
	errs := pr.Errors()
	require.Len(t, errs, 1)
	require.IsType(t, &printer.Error{}, errs[0])
}
