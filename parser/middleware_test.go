package parser

import (
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type VarStatement struct {
	Name  token.Token
	Value ast.Expression
}

func (ls *VarStatement) PrintTo(p *printer.Printer) {
	p.PrintString("var ")
	p.PrintString(ls.Name.Literal)
	p.PrintString(" = ")
	ls.Value.PrintTo(p)
	p.PrintRune(';')
}

func TestUseStatementParser(t *testing.T) {
	l := &lexer.Lexer{}
	l.Init([]byte("var x = 100"))
	p := Parser{}
	p.UseStatementParser(func(p *Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "var" {
			p.AdvanceToken() // consume "var"
			name := p.CurrentToken
			if err := p.Expect(token.IDENT); err != nil {
				return nil, err
			}
			if err := p.Expect(token.ASSIGN); err != nil {
				return nil, err
			}
			value, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			if err := p.ExpectSemi(); err != nil {
				return nil, err
			}
			return &VarStatement{Name: name, Value: value}, nil
		}
		return next()
	})
	p.Init(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatal(err)
	}
	prt := printer.New()
	prog.PrintTo(prt)
	expected := "var x = 100;\n"
	if result := prt.String(); result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
