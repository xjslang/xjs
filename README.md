# XJS (eXtensible JavaScript parser)

[![Go Reference](https://pkg.go.dev/badge/github.com/xjslang/xjs.svg)](https://pkg.go.dev/github.com/xjslang/xjs)
[![Go Report Card](https://goreportcard.com/badge/github.com/xjslang/xjs)](https://goreportcard.com/report/github.com/xjslang/xjs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**XJS** is a highly customizable JavaScript parser written in Go. Our goal is to create a JavaScript compiler that includes only the essential, proven features while enabling users to extend the language through dynamic plugins.

## Key Features

- **🎯 Minimalist Design**: Only essential JavaScript features, no bloat
- **🔧 Extensible Architecture**: Add custom syntax through middlewares
- **⚡ High Performance**: Written in Go for speed and efficiency  

## Quick Start

### Installation

```bash
go get github.com/xjslang/xjs@latest
```

### Basic Usage

```go
package main

import (
	"fmt"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func main() {
	input := `
		let r = 45
		let area = r * r * Math.PI
	`
	l := lexer.New(input)
	p := parser.New(l)
	ast, err := p.ParseProgram()
	if err != nil {
		panic(fmt.Sprintf("ParseProgram() error: %v\n", err))
	}
	fmt.Println(ast.String())
	// Output: let r=45;let area=((r*r)*Math.PI)
}
```

## Philosophy: Minimalism and Sufficiency

Rather than accumulating features over time, **XJS** starts with a carefully curated set of **necessary and sufficient** language constructs. We have deliberately excluded redundant and confusing features:

| Excluded | Reason | Alternative |
|----------|--------|-------------|
| **Classes** | Functions provide sufficient abstraction | Use functions and closures |
| **`const/var`** | Single variable declaration is sufficient | Use `let` only |
| **Weak equality** | `==` automatically becomes `===` | Strict equality only |
| **Arrow functions** | Non-essential | `function() {}` |
| **`try/catch`** | Non-essential | Return error values |
| **`import/export`** | Non-essential | Use `require()` |
| **`async/await`** | Non-essential | Use `.then(onSuccess, onRejected)` |
| **Template literals** | Non-essential | Use `"string " + variable` |
| **Destructuring** | Non-essential | Use dot notation and indexing |
| **Spread operator** | Non-essential | Use `Array.concat()` and `Object.assign()` |
| **Rest parameters** | Non-essential | Use `arguments` keyword |

> [!NOTE]
> You can always create a plugin to implement any excluded features! For example, you might want to create a plugin to support `import/export` statements.

## Extensible Architecture

XJS provides several ways to extend the language:

- **`lexer.Builder.UseInterceptor`** - Add a lexer interceptor
- **`parser.Builder.UseStatementInterceptor`** - Add custom statement types
- **`parser.Builder.UseExpressionInterceptor`** - Add custom expression types  
- **`parser.Builder.UsePrefixOperator`** - Add prefix operators (like `typeof`)
- **`parser.Builder.UseInfixOperator`** - Add infix operators (like `^` for power)
- **`parser.Builder.UseOperand`** - Add custom literals/constants
- **`parser.Builder.UseProgramTransformer`** - Addd a program transformer

### Simple Extension Example

```go
// Add support for the ^ (power) operator
powTokenType := lb.RegisterTokenType("pow")
pb.UseInfixOperator(powTokenType, parser.PRODUCT+1, func(left ast.Expression, right func() ast.Expression) ast.Expression {
    return &PowExpression{Left: left, Right: right()}
})

// Input: let result = 2^3
// Output: let result=Math.pow(2,3)
// Complete example: ./parser/parser_examples_test.go
```

## Documentation

- **[Getting Started Guide](docs/getting-started.md)** - Step-by-step tutorial
- **[Extension Examples](docs/examples/)** - Complete code examples
- **[API Reference](docs/api-reference.md)** - Detailed API documentation
- **[Plugin Development](docs/plugin-development.md)** - Build your own extensions
- **[Architecture Overview](docs/architecture.md)** - Understanding XJS internals

## Ecosystem

Check out these community plugins:

- **[Try-Parser](https://github.com/xjslang/try-parser)** - Adds `try/catch` support
- **[JSX-Parser](https://github.com/xjslang/jsx-parser)** - JSX syntax support
- **[Defer-Parser](https://github.com/xjslang/defer-parser)** - Go-style defer statements

> 🔍 [Explore all plugins](https://github.com/search?q=org%3Axjslang+-parser&type=repositories)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the principles of simplicity and extensibility
- Built with ❤️ by the XJS community
