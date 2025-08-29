package parser

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func TestHandler(t *testing.T) {
	input := `
	let x := "hello" + "Dolly!"`
	l := lexer.New(input)
	p := New(l)
	p.UsePrefixExpressionHandler(token.STRING, func(p *Parser, next Handler) ast.Expression {
		fmt.Println(p.CurrentToken)
		return next(p)
	})
	p.UseInfixExpressionHandler(token.PLUS, func(p *Parser, left ast.Expression, next InfixHandler) ast.Expression {
		fmt.Println(p.CurrentToken)
		return next(p, left)
	})
	p.ParseProgram()
}
