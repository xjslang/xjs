package parser

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func TestHandler(t *testing.T) {
	input := "let x = 'hello' + 'Dolly!' + 100"
	l := lexer.New(input)
	p := New(l)
	p.UseExpressionHandler(func(p *Parser, next func(p *Parser) ast.Expression) ast.Expression {
		switch p.CurrentToken.Type {
		case token.STRING:
			fmt.Println("string!")
		case token.INT:
			fmt.Println("int!")
		}
		return next(p)
	})
	p.ParseProgram()
}
