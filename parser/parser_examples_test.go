package parser

import (
	"fmt"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

type ConstStatement struct {
	Token token.Token
	Name  *ast.Identifier
	Value ast.Expression
}

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
	lb := lexer.NewBuilder()
	pb := NewBuilder(lb)
	// adds support for the `const` keyword!
	pb.UseStatementInterceptor(func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.currentToken.Type == token.IDENT && p.currentToken.Literal == "const" {
			stmt := &ConstStatement{Token: p.currentToken}
			p.NextToken() // moves to identifier token
			stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
			if !p.ExpectToken(token.ASSIGN) { // expects "="
				return nil
			}
			p.NextToken() // moves to value expression
			stmt.Value = p.ParseExpression()
			return stmt
		}
		return next() // otherwise, next!
	})
	p := pb.Build(input)
	ast, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
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
	lb := lexer.NewBuilder()
	// registers typeof keyword
	typeofType := lb.RegisterTokenType("typeof")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.IDENT && ret.Literal == "typeof" {
			ret.Type = typeofType
		}
		return ret
	})
	// adds support for the typeof keyword!
	pb := NewBuilder(lb)
	pb.RegisterPrefixOperator(typeofType, func(token token.Token, right func() ast.Expression) ast.Expression {
		return &TypeofExpression{
			Token: token,
			Right: right(),
		}
	})
	p := pb.Build(input)
	ast, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
	fmt.Println(ast.String())
	// Output: if (((typeof x)==="string")){console.log("x is a string")}
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
	lb := lexer.NewBuilder()
	pb := NewBuilder(lb)
	powType := lb.RegisterTokenType("pow")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		if l.CurrentChar() == '^' {
			l.ReadChar()
			return token.Token{Type: powType, Literal: "^", Line: l.Line(), Column: l.Column()}
		}
		return next()
	})

	// adds support for the ^ operator!
	pb.RegisterInfixOperator(powType, PRODUCT+1, func(token token.Token, left ast.Expression, right func() ast.Expression) ast.Expression {
		return &PowExpression{
			Token: token,
			Left:  left,
			Right: right(),
		}
	})

	p := pb.Build(input)
	ast, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
	fmt.Println(ast.String())
	// Output: let squareArea=Math.pow(r,2)
}

type PiLiteral struct {
	Token token.Token
}

// Tells the parser how to write a node
func (pl *PiLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("Math.PI")
}

func Example_operand() {
	input := "let area = PI * r * r"
	lb := lexer.NewBuilder()
	pb := NewBuilder(lb)
	// registers PI keyword
	piType := lb.RegisterTokenType("PI")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.IDENT && ret.Literal == "PI" {
			ret.Type = piType
		}
		return ret
	})

	// adds support for the PI constant!
	pb.RegisterOperand(piType, func(token token.Token) ast.Expression {
		return &PiLiteral{Token: token}
	})

	p := pb.Build(input)
	ast, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
	fmt.Println(ast.String())
	// Output: let area=((Math.PI*r)*r)
}

type RandomExpression struct {
	Token token.Token
}

func (re *RandomExpression) WriteTo(b *strings.Builder) {
	b.WriteString("Math.random()")
}

func Example_expressionParser() {
	input := "let randomValue = RANDOM + 10"
	lb := lexer.NewBuilder()
	pb := NewBuilder(lb)
	// intercepts expression parsing to handle RANDOM as a special expression!
	pb.UseExpressionInterceptor(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.currentToken.Type == token.IDENT && p.currentToken.Literal == "RANDOM" {
			return p.ParseRemainingExpression(&RandomExpression{Token: p.currentToken})
		}
		return next()
	})
	p := pb.Build(input)
	ast, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
	fmt.Println(ast.String())
	// Output: let randomValue=(Math.random()+10)
}

func Example_combined() {
	input := `
	const circleArea = PI * r^2
	if (typeof radius == 'string') {
		let randomRadius = RANDOM * 10
	}`
	lb := lexer.NewBuilder()
	pb := NewBuilder(lb)
	// registers PI keyword
	piType := lb.RegisterTokenType("PI")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.IDENT && ret.Literal == "PI" {
			ret.Type = piType
		}
		return ret
	})
	// regists infix `^`
	powType := lb.RegisterTokenType("pow")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		if l.CurrentChar() == '^' {
			l.ReadChar() // consume ^
			return token.Token{Type: powType, Literal: "^", Line: l.Line(), Column: l.Column()}
		}
		return next()
	})
	// registers prefix `typeof`
	typeofType := lb.RegisterTokenType("typeof")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.IDENT && ret.Literal == "typeof" {
			ret.Type = typeofType
		}
		return ret
	})

	// combines all previous examples!
	pb.UseStatementInterceptor(func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.currentToken.Type == token.IDENT && p.currentToken.Literal == "const" {
			stmt := &ConstStatement{Token: p.currentToken}
			p.NextToken()
			stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
			if !p.ExpectToken(token.ASSIGN) {
				return nil
			}
			p.NextToken()
			stmt.Value = p.ParseExpression()
			return stmt
		}
		return next()
	})
	pb.RegisterPrefixOperator(typeofType, func(token token.Token, right func() ast.Expression) ast.Expression {
		return &TypeofExpression{Token: token, Right: right()}
	})
	pb.RegisterInfixOperator(powType, PRODUCT+1, func(token token.Token, left ast.Expression, right func() ast.Expression) ast.Expression {
		return &PowExpression{Token: token, Left: left, Right: right()}
	})
	pb.RegisterOperand(piType, func(token token.Token) ast.Expression {
		return &PiLiteral{Token: token}
	})
	pb.UseExpressionInterceptor(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.currentToken.Type == token.IDENT && p.currentToken.Literal == "RANDOM" {
			return p.ParseRemainingExpression(&RandomExpression{Token: p.currentToken})
		}
		return next()
	})
	p := pb.Build(input)
	ast, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
	fmt.Println(ast.String())
	// Output: const circleArea=(Math.PI*Math.pow(r,2));if (((typeof radius)==="string")){let randomRadius=(Math.random()*10)}
}
