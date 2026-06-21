package jsx

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	startTag  = token.RegisterType("start-tag")
	endTag    = token.RegisterType("end-tag")
	concatTag = token.RegisterType("concat")
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

func ParseTag(p *parser.Parser) (_ *Tag, err error) {
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

// Plugin enriches the JavaScript parser, so that we can parse expressions that are not part of the JS standard.
func Plugin(b *builder.Builder) {
	token.RegisterUnaryType(startTag)
	token.RegisterBinaryType(concatTag, token.OR.Precedence())

	// now the parser can "scan" '<' and '</'
	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
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
	})

	// now the parser can "parse" HTML tags
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type == startTag {
			return ParseTag(p)
		}
		return next()
	})

	// now the parser can "concatenate" elements
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (_ ast.Expr, err error) {
		if p.CurrentToken.Type == concatTag {
			node := &ConcatExpr{Left: left}
			p.AdvanceToken()
			if node.Right, err = js.ParseRightExpr(p, concatTag.Precedence()); err != nil {
				return
			}
			return node, nil
		}
		return next(left)
	})
}
