# XJS (eXtensible JavaScript parser)

[![Go Reference](https://pkg.go.dev/badge/github.com/xjslang/xjs.svg)](https://pkg.go.dev/github.com/xjslang/xjs)
[![Go Report Card](https://goreportcard.com/badge/github.com/xjslang/xjs)](https://goreportcard.com/report/github.com/xjslang/xjs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A "pluggable" JavaScript parser.

**XJS** is not a new language; it's just a simplified version of JavaScript without unnecessary, confusing, redundant or deprecated constructs, such as `async/await`, `classes`, `arrow functions`, `with` etc.

The question then becomes: **why would anyone want to use a subset of JavaScript without all those cool features that have accumulated over time?** And there are multiple reasons:

1. Some of those features are questionable, such as `with` which was discouraged by JavaScript itself.
2. Some people prefer a clean language, without unnecessary, confusing, or redundant constructs.
3. Maintaining the base code would be much simpler, as we wouldn't have to maintain those additional and questionable features.

> **But the most important reason is to provide an extensible (or pluggable) parser.**

For example, you might want to use [jsx-parser](https://github.com/xjslang/jsx-parser) to parse JSX syntax, or [defer-parser](https://github.com/xjslang/defer-parser), etc. Installing plugins is very simple (Go):

```go
lb := lexer.NewBuilder()
parser := parser.NewBuilder(lb).
  Install(jsxparser.Plugin).   // adds support for JSX syntax
  Install(deferparser.Plugin). // adds support for "defer" syntax
  // ... more plugins ...
  Build(input)

// the parser now "understands" the JSX and "defer" syntax
// and can translate your custom JavaScript code into standard JavaScript
program, _ := parser.ParseProgram()
fmt.Println(program.String())
```

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
package example

import (
	"fmt"
	"strings"

	interpparser "github.com/xjslang/interp-parser"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func main() {
	input := `
	// interp-parser enhances the language
	// by allowing string interpolation
	console.log('Hello ${name} ${surname}!')`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).
		Install(interpparser.Plugin).
		Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		panic(fmt.Sprintf("ParseProgram() error: %v\n", err))
	}
	fmt.Println(strings.ReplaceAll(program.String(), ";", ";\n"))
	// Output:
	// console.log("Hello ${name} ${surname}!")
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

> You can always create a plugin to implement any excluded features! For example, you might want to create a plugin to support `import/export` statements.

## Extensible Architecture

XJS provides several ways to extend the language:

- **`lexer.Builder.UseTokenInterceptor`** - Intercept and modify token lexing flow
- **`parser.Builder.UseStatementInterceptor`** - Intercept and modify statement parsing flow
- **`parser.Builder.UseExpressionInterceptor`** - Intercept and modify expression parsing flow  
- **`parser.Builder.UsePrefixOperator`** - Register prefix operators (like `typeof`)
- **`parser.Builder.UseInfixOperator`** - Register infix operators (like `^` for power)
- **`parser.Builder.UseOperand`** - Register custom literals/constants
- **`parser.Builder.UseProgramTransformer`** - Transform the generated AST

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

## Ecosystem

Check out these community plugins:

| Plugin | Description |
|--------|-------------|
| **[Defer-Parser](https://github.com/xjslang/defer-parser)** | Adds support for Go-style `defer` statements |
| **[Interp-Parser](https://github.com/xjslang/interp-parser)** | Simplified template literal syntax with variable interpolation |
| **[JSX-Parser](https://github.com/xjslang/jsx-parser)** | JSX syntax support with React.createElement transformation |
| **[Pow-Parser](https://github.com/xjslang/pow-parser)** | Adds support for the `**` power operator |
| **[Spread-Parser](https://github.com/xjslang/spread-parser)** | Object spread syntax using `..` operator (transpiles to Object.assign) |
| **[Switch-Parser](https://github.com/xjslang/switch-parser)** | Adds support for `switch/case` statements |
| **[Try-Parser](https://github.com/xjslang/try-parser)** | Adds `try/catch/finally` error handling support |

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
