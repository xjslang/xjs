package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type LetStatement struct {
	Name  token.Token
	Value ast.Expression
}

func (node *LetStatement) PrintTo(p *printer.Printer) {
	p.PrintString("let ")
	p.PrintString(node.Name.Literal)
	p.PrintString(" = ")
	node.Value.PrintTo(p)
	p.PrintRune(';')
}

func LetPlugin(p *parser.Builder) {
	p.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == token.LET {
			stmt := &LetStatement{}
			p.AdvanceToken() // consume let
			ident := p.CurrentToken
			if err := p.Expect(token.IDENT); err != nil {
				return nil, err
			}
			stmt.Name = ident
			if err := p.Expect(token.ASSIGN); err != nil {
				return nil, err
			}
			val, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			stmt.Value = val
			if err := p.ExpectSemi(); err != nil {
				return nil, err
			}
			return stmt, nil
		}
		return next()
	})
}
