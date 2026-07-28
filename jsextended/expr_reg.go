package jsextended

import (
	"slices"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var validRegexFlags = []byte{'d', 'g', 'i', 'm', 's', 'u', 'v', 'y'}

type RegExpr struct {
	ast.BaseExpr
	Value token.Token
}

func ParseRegExpr(p *parser.Parser) (node *RegExpr, err error) {
	node = &RegExpr{}
	if node.Value, err = p.Expect(token.DIVIDE); err != nil {
		return
	}
	str := strings.Builder{}
	for {
		if p.CurrentToken.AfterNewline {
			break
		}
		if typ := p.CurrentToken.Type; typ == token.DIVIDE || typ == token.EOF {
			break
		}
		if p.CurrentToken.Literal == "\\" && p.PeekToken.Type == token.DIVIDE {
			for range 2 {
				str.Write(p.CurrentToken.Raw)
				p.AdvanceToken()
			}
		}
		str.Write(p.CurrentToken.Raw)
		p.AdvanceToken()
	}
	var close token.Token
	if close, err = p.Expect(token.DIVIDE); err != nil {
		return
	}
	str.Write(close.Raw)
	if !p.CurrentToken.AfterNewline && p.CurrentToken.Type == token.IDENT {
		flags := p.CurrentToken.Raw
		if len(flags) > 0 && flags[0] != ' ' {
			for _, flag := range flags {
				if !slices.Contains(validRegexFlags, flag) {
					err = p.Error("invalid flags")
					return
				}
			}
			str.Write(p.CurrentToken.Raw)
			p.AdvanceToken()
		}
	}
	node.Value.Literal += str.String()
	return
}

func PrintRegExpr(pr *printer.Printer, node *RegExpr) error {
	pr.Print(node.Value)
	return nil
}
