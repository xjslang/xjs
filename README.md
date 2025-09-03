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

This project uses [Mage](https://magefile.org/) as a taks manager:
```bash
go install github.com/magefile/mage@latest

# view available tasks
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
