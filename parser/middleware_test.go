package parser_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type VarStatement struct {
	Name  *ast.Identifier
	Value ast.Expression
}

func (ls *VarStatement) PrintTo(p *printer.Printer) {
	p.PrintString("var ")
	ls.Name.PrintTo(p)
	p.PrintString(" = ")
	ls.Value.PrintTo(p)
	p.PrintRune(';')
}

func TestUseStatementParser(t *testing.T) {
	l := &lexer.Lexer{}
	l.Init([]byte("var x = 100"))
	p := parser.Parser{}
	p.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "var" {
			p.AdvanceToken() // consume "var"

			if p.CurrentToken.Type != token.IDENT {
				msg := fmt.Sprintf("Expected %v", token.IDENT)
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			name := &ast.Identifier{Value: p.CurrentToken.Literal}
			p.AdvanceToken()

			if p.CurrentToken.Type != token.ASSIGN {
				msg := fmt.Sprintf("Expected %v", token.ASSIGN)
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			p.AdvanceToken()

			value, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}

			if p.CurrentToken.Type == token.SEMICOLON {
				p.AdvanceToken()
			} else if p.CurrentToken.Type != token.EOF && !p.CurrentToken.AfterNewline {
				msg := "Expected statement terminator"
				p.AddError(msg)
				return nil, errors.New(msg)
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
