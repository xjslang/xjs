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

// Intercepts the statements and add your own syntax
func ConstStatementHandler(p *Parser, next func() ast.Statement) ast.Statement {
	if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "const" {
		stmt := &ConstStatement{Token: p.CurrentToken}
		// moves to identifier token
		p.NextToken()
		stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
		// expects "="
		if !p.ExpectToken(token.ASSIGN) {
			return nil
		}
		// moves to value and parses it
		p.NextToken()
		stmt.Value = p.ParseExpression()
		return stmt
	}
	// otherwise, next!
	return next()
}

// Represents the PI constant
type PiLiteral struct {
	Token token.Token
}

// Tells the parser how to write a node
func (nl *PiLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("3.1416")
}

// Intercepts the expressions and add your own syntax
func PiExpressionHandler(p *Parser, next func() ast.Expression) ast.Expression {
	if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "PI" {
		return &PiLiteral{Token: p.CurrentToken}
	}
	// otherwise, next!
	return next()
}

// Example_const demonstrates how to create a custom statement parser for the `const` keyword.
func Example_const() {
	input := "const x = 42"

	l := lexer.New(input)
	p := New(l)

	// Register the const statement handler
	p.UseStatementHandler(ConstStatementHandler)

	ast := p.ParseProgram()
	fmt.Println(ast.String())

	// Output: const x=42
}

// Example_pi demonstrates how to create a custom expression parser for the PI constant.
func Example_pi() {
	input := "let pi = PI"

	l := lexer.New(input)
	p := New(l)

	// Register the PI expression handler
	p.UseExpressionHandler(PiExpressionHandler)

	ast := p.ParseProgram()
	fmt.Println(ast.String())

	// Output: let pi=3.1416
}
