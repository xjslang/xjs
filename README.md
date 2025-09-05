# XJS (eXtensible JavaScript parser)

**XJS** is a highly customizable JavaScript parser. The idea is to keep the core minimal, excluding redundant, confusing or non-essential constructs, and allowing users to add their own features.

Check out [VISION.md](./VISION.md) to learn more.

## Supported Language Features

Many JavaScript features are already supported by **XJS**. The following table highlights only the most notable ones:

| Construct                | Supported | Reason                              | Alternative |
|--------------------------|-----------|-------------------------------------|-------------------|
| `let`                    | ✅ Yes    | Essential                           |                   |
| `function`               | ✅ Yes    | Essential                           |                   |
| `if/else`, `while/for`   | ✅ Yes    | Essential                           |                   |
| `==/!=`                  | ✅ Yes    | Essential (translated to `===/!==`) |                   |
| `x++/x--`, `++x/--x`     | ✅ Yes    | Very convenient                     |                   |
| `===/!==`                | ❌ No     | Confusing                           | Use `==/!=`       |
| `const`, `var`           | ❌ No     | Redundant                           | Use `let`         |
| `arrow functions`        | ❌ No     | Redundant                           | Use `function`    |
| `class`                  | ❌ No     | Non-essential                       | Use `function(s)` |
| `switch/case`            | ❌ No     | Non-essential                       | Use `if/else`     |
| `try/catch`              | ❌ No     | Non-essential                       | Use [`try-parser`](https://github.com/xjslang/try-parser) |
| `async/await`            | ❌ No     | Non-essential                       | Create your own parser 😊 |

## Installation

```bash
go get github.com/xjslang/xjs
```

This package uses [Mage](https://magefile.org/) as a taks manager:
```bash
# install mage if not already installed
go install github.com/magefile/mage@latest

# execute the following command from
# the package directory to view available tasks
mage -l
```

## Create your own parser that extend the XJS syntax

Creating your own **XJS** parser is really simple. Just intercept the statements or expressions with the `UseStatementHandler` or `UseExpressionHandler` methods.

### Create a statement parser

**XJS** doesn't support the `const` keyword, but if you are a "const believer", you can create your own plugin. For example:

```go
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
func ConstStatementHandler(p *parser.Parser, next func() ast.Statement) ast.Statement {
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
```

See the complete example [here](./parser/parser_examples_test.go)

### Create an expression parser

```go
// Represents the PI constant
type PiLiteral struct {
	Token token.Token
}

// Tells the parser how to write a node
func (nl *PiLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("3.1416")
}

// Intercepts the expressions and add your own syntax
func PiExpressionHandler(p *parser.Parser, next func() ast.Expression) ast.Expression {
	if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "PI" {
		return &PiLiteral{Token: p.CurrentToken}
	}
	// otherwise, next!
	return next()
}
```

See the complete example [here](./parser/parser_examples_test.go)

### Concatenate multiple statement and expressions parsers:

You can concatenate as many parsers as you want, enriching the parser to your own preferences:

```go
l := lexer.New(input)
p := parser.New(l)

// concatenates multiple parsers that enrich XJS syntax
p.UseStatementHandler(ConstStatementHandler)
p.UseStatementHandler(TryCatchStatementHandler)
p.UseStatementHandler(AwaitStatementHandler)
p.UseExpressionHandler(JsxExpressionHandler)
p.UseExpressionHandler(MathExpressionHandler)
p.UseExpressionHandler(VectorExpressionHandler)

ast := p.ParseProgram()
fmt.Println(ast.String())
```

Here you will find different parsers to inspire you:  
https://github.com/search?q=org%3Axjslang+-parser&type=repositories

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
