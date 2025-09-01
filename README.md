# XJS (eXtensible JavaScript language)

A clean, simplified JavaScript-like language lexer and parser written in Go. This library provides tokenization and parsing capabilities for a JavaScript-like syntax without redundant or complex features like classes, arrow functions, const/var declarations, etc.

## Features

- **Clean JavaScript-like syntax** - Supports essential JavaScript constructs
- **No redundant features** - Excludes classes, arrow functions, multiple variable declaration types
- **Comprehensive parsing** - Full AST generation with position tracking
- **Easy to use** - Simple API for lexing and parsing
- **Well documented** - Extensive documentation and examples
- **Line comments** - Supports `//` line comments for code documentation

## Supported Language Features

### Variables
```javascript
let x = 5  // Integer variable
let name = "Hello World"  // String variable
let pi = 3.14159  // Float variable
```

### Functions
```javascript
// Function declaration
function add(a, b) {
    return a + b  // Return the sum
}
```

### Control Flow
```javascript
// Conditionals
if (x < y) {
    console.log("x is less than y")
} else {
    console.log("x is not less than y")
}

// While loops
while (i < 10) {
    i++  // Postfix increment
}

// For loops
for (let i = 0; i < 10; i++) {  // Both prefix and postfix supported
    console.log(i)
}
```

### Data Types
```javascript
let number = 42
let float = 3.14
let string = "Hello"
let boolean = true
let nothing = null
let array = [1, 2, 3]
let object = {name: "John", age: 30}
```

### Operators
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=` (Note: `==` and `!=` perform strict comparison without type coercion)
- Logical: `&&`, `||`, `!`
- Assignment: `=`
- Increment/Decrement: `++`, `--` (both prefix and postfix)
- Assignment: `=`
- Increment/Decrement: `++`, `--`

## Installation

```bash
go get github.com/xjslang/xjs
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

    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    if len(p.Errors()) > 0 {
        for _, err := range p.Errors() {
            fmt.Println("Error:", err)
        }
        return
    }

    fmt.Println("AST:", program.String())
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
