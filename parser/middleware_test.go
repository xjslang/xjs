package parser

import (
	"fmt"
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
	p.UseStatementParser(func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "var" {
			p.AdvanceToken() // consume "var"

			if p.CurrentToken.Type != token.IDENT {
				p.AddError(fmt.Sprintf("Expected %v", token.IDENT))
				return nil
			}
			name := p.CurrentToken
			p.AdvanceToken()

			if p.CurrentToken.Type != token.ASSIGN {
				p.AddError(fmt.Sprintf("Expected %v", token.ASSIGN))
				return nil
			}
			p.AdvanceToken()

			value, err := p.ParseExpression()
			if err != nil {
				return nil
			}

			if p.CurrentToken.Type == token.SEMICOLON {
				p.AdvanceToken()
			} else if p.CurrentToken.Type != token.EOF && !p.CurrentToken.AfterNewline {
				p.AddError("Expected statement terminator")
				return nil
			}

			return &VarStatement{Name: name, Value: value}
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
