package parser

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func TestHandler(t *testing.T) {
	input := "let x = 'hello' + 'Dolly!'"
	l := lexer.New(input)
	p := New(l)
	p.UsePrefixExpressionHandler(token.STRING, func(p *Parser, next PrefixParseFn) ast.Expression {
		fmt.Println("prefix: ", p.CurrentToken)
		return next(p)
	})
	p.UseInfixExpressionHandler(token.PLUS, func(p *Parser, left ast.Expression, next InfixParseFn) ast.Expression {
		fmt.Println("infix: ", p.CurrentToken)
		return next(p, left)
	})
	p.UseStatementHandler(func(p *Parser, next StatementParseFn) ast.Statement {
		fmt.Println("statement: ", p.CurrentToken)
		return next(p)
	})
	p.ParseProgram()
}
