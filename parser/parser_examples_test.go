package parser

import (
	"fmt"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

// Represents a `const` node
type ConstStatement struct {
	Token token.Token
	Name  *ast.Identifier
	Value ast.Expression
}

// Tells the parser how to write a node
func (ls *ConstStatement) WriteTo(b *strings.Builder) {
	b.WriteString("const ")
	ls.Name.WriteTo(b)
	if ls.Value != nil {
		b.WriteRune('=')
		ls.Value.WriteTo(b)
	}
}

func Example_const() {
	input := "const x = 42"
	l := lexer.New(input)
	p := New(l)
	// adds support for the `const` keyword!
	p.UseStatementParser(func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "const" {
			stmt := &ConstStatement{Token: p.CurrentToken}
			p.NextToken() // moves to identifier token
			stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
			if !p.ExpectToken(token.ASSIGN) { // expects "="
				return nil
			}
			p.NextToken() // moves to value expression
			stmt.Value = p.ParseExpression()
			return stmt
		}
		return next() // otherwise, next!
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: const x=42
}

// Represents a `PI` literal node
type PiLiteral struct {
	Token token.Token
}

// Tells the parser how to write a node
func (pl *PiLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("Math.PI")
}

func Example_pi() {
	input := "let area = PI * r * r"
	l := lexer.New(input)
	p := New(l)
	// Intercepts the expressions and add your own syntax
	p.UseExpressionParser(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "PI" {
			return p.ParseRemainingExpression(&PiLiteral{Token: p.CurrentToken})
		}
		return next()
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: let area=((Math.PI*r)*r)
}
