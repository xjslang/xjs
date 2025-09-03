# XJS (eXtensible JavaScript parser)

**XJS** is a highly customizable JavaScript parser and dialect. The idea is to keep the core minimal, excluding redundant or non essential constructs such as `const`, `var`, or `arrow functions`, and allowing users to add their own constructs through the `UseStatementHandler` and `UseExpressionHandler` methods, which follow a "middleware design pattern" similar to Express.js. Check out [VISION.md](./VISION.md) to learn more.

## Supported Language Features

Many JavaScript features are already supported by **XJS**. The following table highlights only the most notable differences:

| Construct                | Supported | Comments / Alternative            |
|--------------------------|-----------|-----------------------------------|
| `let`                    | ✅ Yes    | Essential                         |
| `function`               | ✅ Yes    | Essential                         |
| `if/else`, `while/for`   | ✅ Yes    | Essential                         |
| `x++/x--`, `++x/--x`     | ✅ Yes    | Very convenient                   |
| `==/!=`                  | ✅ Yes    | Translated to `===/!==`           |
| `===/!==`                | ❌ No     | Use `==/!==`                      |
| `const`, `var`           | ❌ No     | Use `let`                         |
| `arrow functions`        | ❌ No     | Use `function`                    |
| `class`                  | ❌ No     | Use `function(s)`                 |
| `switch/case`            | ❌ No     | Use `if/else`                     |
| `try/catch`              | ❌ No     | Use [`try-parser`](https://github.com/xjslang/try-parser) |
| `async/await`            | ❌ No     | Create your own plugin 😊         |

## Installation

```bash
go get github.com/xjslang/xjs
```

This package uses [Mage](https://magefile.org/) as a taks manager:
```bash
go install github.com/magefile/mage@latest

# execute the following command from
# the package directory to view available tasks
mage -l
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/xjslang/xjs/lexer"
    "github.com/xjslang/xjs/parser"
)

func main() {
    input := `
        let x = 5
        function add(a, b) {
            return a + b
        }
    `

    // translates the input to small tokens
    // and generates the Abstract Syntax Tree (AST)
    l := lexer.New(input)
    p := parser.New(l)
    ast := p.ParseProgram()

    if len(p.Errors()) > 0 {
        for _, err := range p.Errors() {
            fmt.Println("Error:", err)
        }
        return
    }

    // prints AST to JavaScript code
    fmt.Println(ast.String())
}
```

## Create your own parser (or plugin)

Crear un nuevo parser que extienda la sintaxis **XJS** es realmente sencillo.

```go
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

// Intercepts the expressions
func PiExpressionHandler(p *Parser, next func() ast.Expression) ast.Expression {
	if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "PI" {
		return &PiLiteral{Token: p.CurrentToken}
	}
    // otherwise, next!
	return next()
}

func Example_plugin() {
	input := "const pi = PI"
	l := lexer.New(input)
	p := New(l)
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

## API Documentation

### Lexer

The lexer tokenizes input source code:

```go
import "github.com/xjslang/xjs/lexer"

l := lexer.New(sourceCode)
for {
    token := l.NextToken()
    if token.Type == token.EOF {
        break
    }
    fmt.Println(token)
}
```

### Parser

The parser builds an Abstract Syntax Tree (AST):

```go
import "github.com/xjslang/xjs/parser"

l := lexer.New(sourceCode)
p := parser.New(l)
program := p.ParseProgram()

// Check for parsing errors
if len(p.Errors()) > 0 {
    // Handle errors
}
```

### Convenience Function

Use the convenience function for simple parsing:

```go
import "github.com/xjslang/xjs"

program, errors := xjs.Parse(sourceCode)
```

## Examples

See the [examples](examples/) directory for more comprehensive examples.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
