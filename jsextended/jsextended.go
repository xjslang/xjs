package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func Plugin(b *plugin.Builder) {
	token.RegisterUnaryType(NEW)

	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		switch tok.Type {
		case token.IDENT:
			switch tok.Literal {
			case "new":
				tok.Type = NEW
			}
		}
		return
	})
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case NEW:
			return ParseNewExpr(p)
		}
		return next()
	})
}

func Printer(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *NewExpr:
		return PrintNewExpr(pr, v)
	}
	return next(node)
}
