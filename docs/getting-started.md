# Getting Started with XJS

Welcome to XJS! This guide will help you get up and running with the XJS parser in just a few minutes.

## Installation

XJS requires Go 1.21 or later. Install it using:

```bash
go get github.com/xjslang/xjs@latest
```

## Your First Parser

Let's start with a simple example:

```go
package main

import (
	"fmt"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func main() {
	// JavaScript code to parse
	input := `
		let name = "XJS"
		let version = 1.0
		console.log("Hello from " + name + " v" + version)
	`
	
	// Create lexer and parser
	l := lexer.New(input)
	p := parser.New(l)
	
	// Parse the code
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}
	
	// Check for parsing errors
	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err.Message)
		}
		return
	}
	
	// Print the parsed result
	fmt.Println("Parsed successfully!")
	fmt.Println("Output:", program.String())
}
```

## Understanding the Output

XJS transforms the input JavaScript into a normalized form:

- **Variables**: Only `let` declarations are used
- **Equality**: `==` becomes `===` automatically
- **Formatting**: Consistent spacing and parentheses

Example transformations:

| Input | Output |
|-------|--------|
| `let x = 5` | `let x=5` |
| `if (a == b)` | `if ((a===b))` |
| `x + y * z` | `(x+(y*z))` |

## What's Supported

XJS supports a curated set of JavaScript features:

### ✅ Supported Features
- **Variables**: `let` declarations
- **Functions**: Function declarations and expressions
- **Control Flow**: `if/else`, `for`, `while`
- **Data Types**: Numbers, strings, booleans, arrays, objects
- **Operators**: Arithmetic, comparison, logical
- **Function Calls**: Including method calls

### ❌ Intentionally Excluded
- **Classes**: Use functions instead
- **Arrow functions**: Use regular functions
- **`const/var`**: Use `let` only
- **`try/catch`**: Use error return values
- **Weak equality**: `==` automatically becomes `===`

## Error Handling

XJS provides detailed error information:

```go
// Example with syntax error
input := `let x = 5 +` // Missing right operand

l := lexer.New(input)
p := parser.New(l)
program, err := p.ParseProgram()

if err != nil {
	fmt.Printf("Parse error: %v\n", err)
	// Output: Parse error: parse error at line 1, column 11: unexpected token EOF
}

// Multiple errors are collected
if len(p.Errors()) > 0 {
	for _, err := range p.Errors() {
		fmt.Printf("Line %d, Col %d: %s\n", 
			err.Position.Line, 
			err.Position.Column, 
			err.Message)
	}
}
```

## Next Steps

Now that you have XJS running, explore these topics:

1. **[Basic Examples](examples/basic-usage.md)** - More parsing examples
2. **[Extension System](examples/middleware-system.md)** - Adding custom features
3. **[API Reference](api-reference.md)** - Complete API documentation
4. **[Plugin Development](plugin-development.md)** - Creating your own extensions

## Common Patterns

### Parsing Files

```go
import (
	"os"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func parseFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	l := lexer.New(string(content))
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return err
	}
	
	fmt.Println(program.String())
	return nil
}
```

### Working with the AST

```go
import "github.com/xjslang/xjs/ast"

// After parsing...
for _, stmt := range program.Statements {
	switch s := stmt.(type) {
	case *ast.LetStatement:
		fmt.Printf("Variable: %s\n", s.Name.Value)
	case *ast.FunctionDeclaration:
		fmt.Printf("Function: %s\n", s.Name.Value)
	case *ast.ExpressionStatement:
		fmt.Println("Expression statement")
	}
}
```

## Troubleshooting

### Common Issues

**Problem**: `unexpected token` errors
**Solution**: Check for missing semicolons or unmatched brackets

**Problem**: Weak equality (`==`) not working
**Solution**: This is intentional! XJS converts `==` to `===` automatically

**Problem**: `const` or `var` not recognized
**Solution**: XJS only supports `let`. Use `let` instead

**Problem**: Classes not working
**Solution**: XJS doesn't support classes. Use functions and closures

### Getting Help

- **GitHub Issues**: Report bugs or ask questions
- **Discussions**: Share ideas and get community help
- **Examples**: Check the [examples directory](examples/) for more code samples

---

Ready to extend XJS with custom features? Check out the [Extension Examples](examples/middleware-system.md)!