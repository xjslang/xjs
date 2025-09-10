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

func Example_statement() {
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

type TypeofExpression struct {
	Token token.Token
	Right ast.Expression
}

func (te *TypeofExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(typeof ")
	te.Right.WriteTo(b)
	b.WriteRune(')')
}

func Example_prefixOperator() {
	input := "if (typeof x == 'string') { console.log('x is a string') }"
	l := lexer.New(input)
	p := New(l)
	// adds support for the typeof keyword!
	p.RegisterPrefixOperator("typeof", func(right func() ast.Expression) ast.Expression {
		return &TypeofExpression{
			Token: p.CurrentToken,
			Right: right(),
		}
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: if ((typeof (x==="string"))){console.log("x is a string")}
}

type PowExpression struct {
	Token token.Token
	Left  ast.Expression
	Right ast.Expression
}

func (pe *PowExpression) WriteTo(b *strings.Builder) {
	b.WriteString("Math.pow(")
	pe.Left.WriteTo(b)
	b.WriteRune(',')
	pe.Right.WriteTo(b)
	b.WriteRune(')')
}

func Example_infixOperator() {
	input := "let squareArea = r^2"
	l := lexer.New(input)
	p := New(l)
	// adds support for the ^ operator!
	p.RegisterInfixOperator("^", PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
		return &PowExpression{
			Token: p.CurrentToken,
			Left:  left,
			Right: right(),
		}
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: let squareArea=Math.pow(r,2)
}

// Represents a `PI` literal node
type PiLiteral struct {
	Token token.Token
}

// Tells the parser how to write a node
func (pl *PiLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("Math.PI")
}

func Example_operand() {
	input := "let area = PI * r * r"
	l := lexer.New(input)
	p := New(l)
	// adds support for the PI constant!
	p.RegisterOperand("PI", func() ast.Expression {
		return &PiLiteral{Token: p.CurrentToken}
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: let area=((Math.PI*r)*r)
}
