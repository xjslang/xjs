package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func NewBuilder() *builder.Builder {
	return builder.New().Install(jsPlugin)
}

func NewPrinter() *printer.Printer {
	p := &printer.Printer{}
	p.UsePrinter(func(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
		switch v := node.(type) {
		case *js.Program:
			js.PrintProgram(p, v)
			return
		case *js.BlockStmt:
			js.PrintBlockStmt(p, v)
			return
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
		case *js.UnaryExpr:
			js.PrintUnaryExpr(p, v)
			return
		case *js.BinaryExpr:
			js.PrintBinaryExpr(p, v)
			return
		case *js.Ident:
			js.PrintIdent(p, v)
			return
		case *js.BasicLit:
			js.PrintBasicLit(p, v)
			return
		case *js.Stmt:
			js.PrintStmt(p, v)
			return
		}
		next(node)
	})
	p.Init()
	return p
}

func Parse(input []byte) (*js.Program, error) {
	p := NewBuilder().Build(input)
	return js.ParseProgram(p)
}

func jsPlugin(b *builder.Builder) {
	token.RegisterUnaryType(token.LPAREN)     //	to evaluate ParenExpr
	token.RegisterBinaryType(token.LPAREN, 7) // to evaluate CallExpr

	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) (tok token.Token) {
		tok = next()
		if tok.Type != token.IDENT {
			return
		}
		switch tok.Literal {
		case "function":
			tok.Type = js.FUNCTION
		case "let":
			tok.Type = js.LET
		}
		return
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		switch p.CurrentToken.Type {
		case js.FUNCTION:
			return js.ParseFunctionDecl(p)
		case js.LET:
			return js.ParseLetStmt(p)
		}
		return js.ParseStmt(p)
	})
	b.UseExprParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		return js.ParseExpr(p)
	})
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == token.LPAREN {
			return js.ParseParenExpr(p)
		}
		return js.ParseUnaryExpr(p)
	})
	b.UseBinaryParser(func(p *parser.Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (ast.Node, error) {
		if p.CurrentToken.Type == token.LPAREN {
			return js.ParseCallExpr(p, leftVal)
		}
		return js.ParseBinaryExpr(p, leftVal)
	})
}
