# Middleware System Examples

XJS's extensibility comes from its powerful middleware system. This guide shows you how to extend the language with custom features.

## Core Middleware Functions

XJS provides these main extension points:

- **`UseStatementParser`** - Add custom statement types
- **`UseExpressionParser`** - Add custom expression types
- **`RegisterPrefixOperator`** - Add prefix operators
- **`RegisterInfixOperator`** - Add infix operators  
- **`RegisterOperand`** - Add custom literals/constants

## UseStatementParser Examples

### Adding `const` Support

```go
package main

import (
	"fmt"
	"strings"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// Define the AST node for const statements
type ConstStatement struct {
	Token token.Token
	Name  *ast.Identifier
	Value ast.Expression
}

// Implement the WriteTo method for code generation
func (cs *ConstStatement) WriteTo(b *strings.Builder) {
	b.WriteString("const ")
	cs.Name.WriteTo(b)
	if cs.Value != nil {
		b.WriteRune('=')
		cs.Value.WriteTo(b)
	}
}

func main() {
	input := "const PI = 3.14159"
	
	l := lexer.New(input)
	p := parser.New(l)
	
	// Add const statement support
	p.UseStatementParser(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "const" {
			stmt := &ConstStatement{Token: p.CurrentToken}
			p.NextToken() // move to identifier
			stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
			
			if !p.ExpectToken(token.ASSIGN) {
				return nil
			}
			p.NextToken() // move to value
			stmt.Value = p.ParseExpression()
			return stmt
		}
		return next() // fallback to default parsing
	})
	
	program, _ := p.ParseProgram()
	fmt.Println(program.String())
	// Output: const PI=3.14159
}
```

### Adding `import` Statements

```go
type ImportStatement struct {
	Token      token.Token
	ImportPath ast.Expression
	Specifiers []*ast.Identifier
}

func (is *ImportStatement) WriteTo(b *strings.Builder) {
	b.WriteString("import ")
	if len(is.Specifiers) > 0 {
		b.WriteRune('{')
		for i, spec := range is.Specifiers {
			if i > 0 {
				b.WriteRune(',')
			}
			spec.WriteTo(b)
		}
		b.WriteString("} from ")
	}
	is.ImportPath.WriteTo(b)
}

func addImportSupport(p *parser.Parser) {
	p.UseStatementParser(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "import" {
			stmt := &ImportStatement{Token: p.CurrentToken}
			p.NextToken()
			
			// Parse import specifiers: import {a, b} from "module"
			if p.CurrentToken.Type == token.LBRACE {
				p.NextToken()
				for p.CurrentToken.Type != token.RBRACE && p.CurrentToken.Type != token.EOF {
					if p.CurrentToken.Type == token.IDENT {
						stmt.Specifiers = append(stmt.Specifiers, &ast.Identifier{
							Token: p.CurrentToken,
							Value: p.CurrentToken.Literal,
						})
					}
					p.NextToken()
					if p.CurrentToken.Type == token.COMMA {
						p.NextToken()
					}
				}
				p.NextToken() // skip closing brace
				
				// Expect "from"
				if p.CurrentToken.Literal == "from" {
					p.NextToken()
				}
			}
			
			// Parse the module path
			stmt.ImportPath = p.ParseExpression()
			return stmt
		}
		return next()
	})
}
```

## UseExpressionParser Examples

### Template Literals

```go
type TemplateLiteral struct {
	Token        token.Token
	Parts        []ast.Expression // string parts and expressions
	Expressions  []ast.Expression // ${} expressions
}

func (tl *TemplateLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("(")
	for i, part := range tl.Parts {
		if i > 0 {
			b.WriteString("+")
		}
		part.WriteTo(b)
	}
	b.WriteString(")")
}

func addTemplateLiterals(p *parser.Parser) {
	p.UseExpressionParser(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type == token.BACKTICK {
			// Parse template literal: `Hello ${name}!`
			tl := &TemplateLiteral{Token: p.CurrentToken}
			// Implementation would parse the template string
			// and extract expressions from ${} blocks
			return tl
		}
		return next()
	})
}
```

### Async/Await Support

```go
type AsyncExpression struct {
	Token token.Token
	Body  ast.Expression
}

func (ae *AsyncExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(async ")
	ae.Body.WriteTo(b)
	b.WriteString(")")
}

type AwaitExpression struct {
	Token token.Token
	Right ast.Expression
}

func (ae *AwaitExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(await ")
	ae.Right.WriteTo(b)
	b.WriteString(")")
}

func addAsyncAwait(p *parser.Parser) {
	// Add async function support
	p.UseExpressionParser(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "async" {
			asyncExpr := &AsyncExpression{Token: p.CurrentToken}
			p.NextToken()
			asyncExpr.Body = p.ParseExpression()
			return asyncExpr
		}
		return next()
	})
	
	// Add await support
	p.RegisterPrefixOperator("await", func(right func() ast.Expression) ast.Expression {
		return &AwaitExpression{
			Token: p.CurrentToken,
			Right: right(),
		}
	})
}
```

## RegisterPrefixOperator Examples

### typeof Operator

```go
type TypeofExpression struct {
	Token token.Token
	Right ast.Expression
}

func (te *TypeofExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(typeof ")
	te.Right.WriteTo(b)
	b.WriteString(")")
}

func main() {
	input := `if (typeof name === "string") { console.log("It's a string!") }`
	
	l := lexer.New(input)
	p := parser.New(l)
	
	p.RegisterPrefixOperator("typeof", func(right func() ast.Expression) ast.Expression {
		return &TypeofExpression{
			Token: p.CurrentToken,
			Right: right(),
		}
	})
	
	program, _ := p.ParseProgram()
	fmt.Println(program.String())
	// Output: if (((typeof name)==="string")){console.log("It's a string!")}
}
```

### delete Operator

```go
type DeleteExpression struct {
	Token token.Token
	Right ast.Expression
}

func (de *DeleteExpression) WriteTo(b *strings.Builder) {
	b.WriteString("(delete ")
	de.Right.WriteTo(b)
	b.WriteString(")")
}

p.RegisterPrefixOperator("delete", func(right func() ast.Expression) ast.Expression {
	return &DeleteExpression{
		Token: p.CurrentToken,
		Right: right(),
	}
})
```

## RegisterInfixOperator Examples

### Power Operator (^)

```go
type PowerExpression struct {
	Token token.Token
	Left  ast.Expression
	Right ast.Expression
}

func (pe *PowerExpression) WriteTo(b *strings.Builder) {
	b.WriteString("Math.pow(")
	pe.Left.WriteTo(b)
	b.WriteString(",")
	pe.Right.WriteTo(b)
	b.WriteString(")")
}

func main() {
	input := "let result = base^exponent"
	
	l := lexer.New(input)
	p := parser.New(l)
	
	// Register ^ operator with high precedence (higher than multiplication)
	p.RegisterInfixOperator("^", parser.PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
		return &PowerExpression{
			Token: p.CurrentToken,
			Left:  left,
			Right: right(),
		}
	})
	
	program, _ := p.ParseProgram()
	fmt.Println(program.String())
	// Output: let result=Math.pow(base,exponent)
}
```

### Pipeline Operator (|>)

```go
type PipelineExpression struct {
	Token token.Token
	Left  ast.Expression
	Right ast.Expression
}

func (pe *PipelineExpression) WriteTo(b *strings.Builder) {
	// Transform a |> b into b(a)
	pe.Right.WriteTo(b)
	b.WriteString("(")
	pe.Left.WriteTo(b)
	b.WriteString(")")
}

p.RegisterInfixOperator("|>", parser.LOWEST+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
	return &PipelineExpression{
		Token: p.CurrentToken,
		Left:  left,
		Right: right(),
	}
})

// Usage: value |> transform |> display
// Becomes: display(transform(value))
```

## RegisterOperand Examples

### Mathematical Constants

```go
type MathConstant struct {
	Token token.Token
	Value string
}

func (mc *MathConstant) WriteTo(b *strings.Builder) {
	b.WriteString("Math.")
	b.WriteString(mc.Value)
}

func addMathConstants(p *parser.Parser) {
	constants := map[string]string{
		"PI":  "PI",
		"E":   "E",
		"TAU": "PI*2", // τ = 2π
	}
	
	for name, value := range constants {
		p.RegisterOperand(name, func() ast.Expression {
			return &MathConstant{
				Token: p.CurrentToken,
				Value: value,
			}
		})
	}
}

// Usage: let area = PI * radius^2
// Output: let area=(Math.PI*Math.pow(radius,2))
```

### Special Values

```go
type SpecialValue struct {
	Token token.Token
	Name  string
}

func (sv *SpecialValue) WriteTo(b *strings.Builder) {
	switch sv.Name {
	case "RANDOM":
		b.WriteString("Math.random()")
	case "NOW":
		b.WriteString("Date.now()")
	case "INFINITY":
		b.WriteString("Infinity")
	}
}

p.RegisterOperand("RANDOM", func() ast.Expression {
	return &SpecialValue{Token: p.CurrentToken, Name: "RANDOM"}
})

p.RegisterOperand("NOW", func() ast.Expression {
	return &SpecialValue{Token: p.CurrentToken, Name: "NOW"}
})
```

## Combining Multiple Extensions

```go
package main

import (
	"fmt"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func setupEnhancedParser() *parser.Parser {
	l := lexer.New("")
	p := parser.New(l)
	
	// Add all our custom features
	addConstSupport(p)
	addTemplateLiterals(p)
	addAsyncAwait(p)
	addMathConstants(p)
	
	// Add typeof operator
	p.RegisterPrefixOperator("typeof", func(right func() ast.Expression) ast.Expression {
		return &TypeofExpression{Token: p.CurrentToken, Right: right()}
	})
	
	// Add power operator
	p.RegisterInfixOperator("^", parser.PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
		return &PowerExpression{Token: p.CurrentToken, Left: left, Right: right()}
	})
	
	return p
}

func main() {
	input := `
		const radius = 5
		const area = PI * radius^2
		if (typeof area === "number") {
			console.log("Area:", area)
		}
	`
	
	l := lexer.New(input)
	p := setupEnhancedParser()
	// Reset parser with new lexer
	// (In actual implementation, you'd need a Reset method)
	
	program, _ := p.ParseProgram()
	fmt.Println(program.String())
}
```

## Middleware Execution Order

Middlewares are executed in **LIFO (Last-In, First-Out)** order:

```go
p.UseExpressionParser(middleware1) // Executed third
p.UseExpressionParser(middleware2) // Executed second  
p.UseExpressionParser(middleware3) // Executed first

// Execution flow: middleware3 -> middleware2 -> middleware1 -> default
```

This allows you to:
1. **Override** previous middlewares by not calling `next()`
2. **Transform** results from previous middlewares
3. **Chain** multiple transformations

## Best Practices

### 1. Always Implement WriteTo

Every custom AST node must implement the `WriteTo` method:

```go
func (n *MyNode) WriteTo(b *strings.Builder) {
	// Generate the output JavaScript code
	b.WriteString("/* my custom syntax */")
}
```

### 2. Handle Edge Cases

```go
p.UseStatementParser(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
	if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "myKeyword" {
		// Check for unexpected EOF
		if p.CurrentToken.Type == token.EOF {
			return nil
		}
		
		// Parse your custom statement
		return parseMyStatement(p)
	}
	return next()
})
```

### 3. Preserve Token Information

Always store the token for error reporting and debugging:

```go
type MyNode struct {
	Token token.Token // Always include this
	// ... other fields
}
```

### 4. Consider Operator Precedence

When adding infix operators, choose precedence carefully:

```go
// Precedence levels (from lowest to highest)
const (
	LOWEST = iota      // 1
	ASSIGNMENT         // 2: =
	LOGICAL_OR         // 3: ||  
	LOGICAL_AND        // 4: &&
	EQUALITY           // 5: === !==
	COMPARISON         // 6: > < >= <=
	SUM                // 7: + -
	PRODUCT            // 8: * / %
	UNARY              // 9: -x !x
	POSTFIX            // 10: x++ x--
	CALL               // 11: func()
	MEMBER             // 12: obj.prop
)
```

---

Next: Learn about [Advanced Plugin Development](../plugin-development.md)