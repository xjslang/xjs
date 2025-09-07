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
	// Intercepts the statements and add your own syntax
	p.UseStatementHandler(func(p *Parser, next func() ast.Statement) ast.Statement {
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
	p.UseExpressionHandler(func(p *Parser, next func(left ast.Expression) ast.Expression) ast.Expression {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "PI" {
			// Continue parsing the rest of the expression
			return next(&PiLiteral{Token: p.CurrentToken})
		}
		return next(nil)
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: let area=((Math.PI*r)*r)
}
