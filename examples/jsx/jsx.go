package jsx

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	startTag  = token.RegisterType("START_TAG", "<")
	endTag    = token.RegisterType("END_TAG", "</")
	concatTag = token.RegisterType("CONCAT", "|")
)

type Tag struct {
	ast.BaseExpr
	Name     *js.Ident
	Children ast.Expr
}

type ConcatExpr struct {
	ast.BaseExpr
	Left  ast.Expr
	Right ast.Expr
}

func parseTag(p *parser.Parser) (_ *Tag, err error) {
	node := &Tag{}
	if _, err = p.Expect(startTag); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	if _, err = p.Expect(token.GT); err != nil {
		return
	}
	if p.CurrentToken.Type != endTag {
		if node.Children, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if _, err = p.Expect(endTag); err != nil {
		return
	}
	var ident *js.Ident
	if ident, err = js.ParseIdent(p); err != nil {
		return
	}
	if ident.Literal != node.Name.Literal {
		return nil, p.ErrorAt(
			ident.Token,
			"expected closing tag </"+node.Name.Literal+">",
		)
	}
	if _, err = p.Expect(token.GT); err != nil {
		return
	}
	return node, nil
}

func jsxScanner(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
	if tok, err = next(); err != nil {
		return
	}
	if tok.Type == token.LT {
		c := sc.CurrentChar()
		switch {
		case scanner.IsLetter(c):
			tok.Type = startTag
		case c == '/':
			sc.AdvanceChar()
			tok.Type = endTag
			tok.Literal = "</"
		}
	} else if tok.Literal == "|" {
		tok.Type = concatTag
	}
	return
}

func jsxUnaryParser(p *parser.Parser, next func() (ast.Expr, error)) (_ ast.Expr, err error) {
	if p.CurrentToken.Type == startTag {
		return parseTag(p)
	}
	return next()
}

func jsxBinaryParser(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (_ ast.Expr, err error) {
	if p.CurrentToken.Type == concatTag {
		node := &ConcatExpr{Left: left}
		p.AdvanceToken()
		if node.Right, err = js.ParseRightExpr(p, concatTag.Precedence()); err != nil {
			return
		}
		return node, nil
	}
	return next(left)
}

func Parse(input []byte) (*js.Program, error) {
	token.RegisterUnaryType(startTag)
	token.RegisterBinaryType(concatTag, token.OR.Precedence())

	sb := xjs.ScannerBuilder()
	sb.UseScanner(jsxScanner)

	pb := xjs.ParserBuilder()
	pb.UseUnaryParser(jsxUnaryParser)
	pb.UseBinaryParser(jsxBinaryParser)

	sc := sb.Build(input)
	p := pb.Build(sc)
	return js.ParseProgram(p)
}

func Compile(result ast.Node) (string, error) {
	pr := xjs.PrinterBuilder().
		UsePrinter(Compiler).
		Build(printer.Compact())
	pr.Print(result)
	return pr.Output()
}

func Format(result ast.Node, opts ...printer.Option) (string, error) {
	pr := xjs.PrinterBuilder().
		UsePrinter(Formatter).
		Build(opts...)
	pr.Print(result)
	return pr.Output()
}
