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
	let x := 100`
	l := lexer.New(input)
	p := New(l)
	p.UsePrefixHandler(token.INT, func(p *Parser, next Handler) ast.Expression {
		fmt.Println("yes!")
		return next(p)
	})
	p.ParseProgram()
}
