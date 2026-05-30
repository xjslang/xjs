package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func NewScanner() *scanner.Scanner {
	token.LPAREN.RegisterPrefixOp()
	token.LPAREN.RegisterInfixOp(7)
	s := &scanner.Scanner{}
	s.UseScanner(func(sc *scanner.Scanner, next func() token.Token) (tok token.Token) {
		tok = next()
		if tok.Type == token.IDENT {
			switch tok.Literal {
			case "function":
				tok.Type = js.FUNCTION
			case "let":
				tok.Type = js.LET
			}
		}
		return tok
	})
	return s
}

func NewParser() *parser.Parser {
	p := &parser.Parser{}
	p.UsePrefixOpParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == token.LPAREN {
			return js.ParseParenExpr(p)
		}
		return next()
	})
	p.UseInfixOpParser(func(p *parser.Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == token.LPAREN {
			return js.ParseCallExpr(p, leftVal)
		}
		return next(leftVal)
	})
	p.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		switch p.CurrentToken.Type {
		case js.FUNCTION:
			return js.ParseFunctionDecl(p)
		case js.LET:
			return js.ParseLetStmt(p)
		}
		return next()
	})
	return p
}

func NewPrinter() *printer.Printer {
	p := &printer.Printer{}
	p.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		switch v := node.(type) {
		case *js.FunctionDecl:
			js.PrintFunctionDecl(p, v)
			return
		case *js.LetStmt:
			js.PrintLetStmt(p, v)
			return
		case *js.CallExpr:
			js.PrintCallExpr(p, v)
			return
		case *js.ParenExpr:
			js.PrintParenExpr(p, v)
			return
		}
		next(node)
	})
	return p
}
