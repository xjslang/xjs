# XJS (eXtensible JavaScript parser)

**XJS** is a highly customizable JavaScript parser. The idea is to keep the core minimal, excluding redundant or non-essential constructs such as `const`, `var`, or `arrow functions`, and allowing users to add their own constructs through the `UseStatementHandler` and `UseExpressionHandler` methods, which follow a "middleware design pattern" similar to Express.js.

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

## Create your own parser to extend the XJS syntax

Creating a new parser that extends the **XJS** syntax is very simple. You just need to declare the structures you want to add to the language and intercept statements or expressions. In the following example, we have added the `const` statement and the `PI` constant:

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

type ConstStatement struct {
	Token token.Token
	Name  *ast.Identifier
	Value ast.Expression
}

// WriteTo instructs the parser how to write a node
func (ls *ConstStatement) WriteTo(b *strings.Builder) {
	b.WriteString("const ")
	ls.Name.WriteTo(b)
	if ls.Value != nil {
		b.WriteRune('=')
		ls.Value.WriteTo(b)
	}
}

type PiLiteral struct {
	Token token.Token
}

func (nl *PiLiteral) WriteTo(b *strings.Builder) {
	b.WriteString("3.1416")
}

// Intercepts the statements
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

// Intercepts the expressions
func PiExpressionHandler(p *parser.Parser, next func() ast.Expression) ast.Expression {
	if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "PI" {
		return &PiLiteral{Token: p.CurrentToken}
	}
	// otherwise, next!
	return next()
}

func main() {
	input := "const pi = PI"
	l := lexer.New(input)
	p := parser.New(l)

	// extends the language syntax!
	p.UseStatementHandler(ConstStatementHandler)
	p.UseExpressionHandler(PiExpressionHandler)

	ast := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Println("Error:", err)
		}
		return
	}
	fmt.Println(ast.String())
	// Output: const pi=3.1416
}
```

Additionally, we can chain as many parsers as we want:

```go
// ...
p := New(l)
p.UseStatementHandler(ConstStatementHandler)
p.UseStatementHandler(TryCatchStatementHandler)
p.UseStatementHandler(AwaitStatementHandler)
p.UseExpressionHandler(JsxExpressionHandler)
p.UseExpressionHandler(MathExpressionHandler)
p.UseExpressionHandler(VectorExpressionHandler)
// ...
```

Here you will find numerous parser examples:  
https://github.com/search?q=org%3Axjslang+-parser&type=repositories

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
