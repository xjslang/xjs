# XJS (eXtensible JavaScript parser)

**XJS** is a highly customizable JavaScript parser. Our goal is to create a JavaScript compiler that includes only the essential, proven features while enabling users to extend the language through dynamic plugins.

## Minimalism and Sufficiency

Rather than accumulating features over time, **XJS** starts with a carefully curated set of **necessary and sufficient** language constructs. We have deliberately excluded redundant features:

- **No classes** - Functions provide sufficient abstraction capabilities
- **No arrow functions** - Regular function syntax is adequate
- **No `const/var`** - A single variable declaration mechanism suffices
- **No `try/catch`** - Alternative error handling patterns are preferred
- **No redundant syntactic sugar** - Focus on core functionality

This approach ensures that every included feature has demonstrated genuine utility and necessity over the years.

## Installation

```bash
go get github.com/xjslang/xjs@latest
```

## Create your own parser that extends the XJS syntax

Creating your own **XJS** parser is really simple. Just intercept the statements or expressions with the `UseStatementParser`, `UseExpressionParser` and `UseRemainingExpressionParser` methods.

<details>
	<summary>Example 1: Create a statement parser</summary>

```go
// XJS doesn't support the `const` keyword, but if you are a "const believer", you can create your own plugin.

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
	p.UseStatementParser(func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "const" {
			stmt := &ConstStatement{Token: p.CurrentToken}
			p.NextToken() // moves to identifier token
			stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
			if !p.ExpectToken(token.ASSIGN) { // expects "="
				return nil
			}
			p.NextToken() // moves to value and parses it
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
```

</details>

<details>
	<summary>Example 2: Create an expression parser</summary>

```go
// In the following example we are going to declare the `PI` literal

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
			// uses the new node and continues parsing the rest of the expression
			return p.ParseRemainingExpression(&PiLiteral{Token: p.CurrentToken})
		}
		return next()
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: let area=((Math.PI*r)*r)
}
```

</details>

<details>
	<summary>Example 3: Create an operator parser</summary>

```go
// In the following example we are going to create a mathematical pow example:

// Represents a Power expression
type PowExpression struct {
	Token token.Token
	Left  ast.Expression
	Right ast.Expression
}

// Tells the parser how to write a node
func (pe *PowExpression) WriteTo(b *strings.Builder) {
	b.WriteString("Math.pow(")
	pe.Left.WriteTo(b)
	b.WriteRune(',')
	pe.Right.WriteTo(b)
	b.WriteRune(')')
}

func Example_pow() {
	input := "let y = x + r ^ r"
	l := lexer.New(input)
	p := New(l)
	p.UseRemainingExpressionParser(func(p *Parser, left ast.Expression, next func() ast.Expression) ast.Expression {
		if p.PeekToken.Type == token.ILLEGAL && p.PeekToken.Literal == "^" {
			p.NextToken() // consume the ^ token
			p.NextToken() // move to the right operand
			exp := &PowExpression{
				Token: p.CurrentToken,
				Left:  left,
				Right: p.ParseExpression(),
			}
			return exp
		}
		return next()
	})
	ast := p.ParseProgram()
	fmt.Println(ast.String())
	// Output: let y=(x+Math.pow(r,r))
}
```

</details>

<details>
	<summary>Example 4: Concatenate multiple parsers</summary>

```go
// You can concatenate as many parsers as you want, enriching the parser to your own preferences:

l := lexer.New(input)
p := parser.New(l)

// concatenates multiple parsers that enrich XJS syntax
p.UseStatementParser(ConstStatementParser)
p.UseStatementParser(TryCatchStatementParser)
p.UseStatementParser(AwaitStatementParser)
p.UseExpressionParser(JsxExpressionParser)
p.UseExpressionParser(MathExpressionParser)
p.UseExpressionParser(VectorExpressionParser)
p.UseRemainingExpressionParser(PowExpressionParser)
p.UseRemainingExpressionParser(XORExpressionParser)

ast := p.ParseProgram()
fmt.Println(ast.String())
```

</details>

Here you will find different parsers to inspire you:  
https://github.com/search?q=org%3Axjslang+-parser&type=repositories

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
