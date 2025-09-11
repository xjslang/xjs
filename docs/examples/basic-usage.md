# Basic Usage Examples

This document contains practical examples of using XJS for common tasks.

## Simple Parsing

### Example 1: Basic Variable Assignment

```go
package main

import (
	"fmt"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func main() {
	input := `let message = "Hello, World!"`
	
	l := lexer.New(input)
	p := parser.New(l)
	program, _ := p.ParseProgram()
	
	fmt.Println(program.String())
	// Output: let message="Hello, World!"
}
```

### Example 2: Function Declaration

```go
input := `
function greet(name) {
	return "Hello, " + name + "!"
}
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
// Output: function greet(name){return ("Hello, "+(name+"!"))}
```

### Example 3: Control Flow

```go
input := `
let age = 25
if (age >= 18) {
	console.log("Adult")
} else {
	console.log("Minor")
}
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
// Output: let age=25;if ((age>=18)){console.log("Adult")}else{console.log("Minor")}
```

## Working with Data Structures

### Arrays

```go
input := `
let numbers = [1, 2, 3, 4, 5]
let fruits = ["apple", "banana", "orange"]
numbers.push(6)
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
// Output: let numbers=[1,2,3,4,5];let fruits=["apple","banana","orange"];numbers.push(6)
```

### Objects

```go
input := `
let person = {
	name: "John",
	age: 30,
	city: "New York"
}
console.log(person.name)
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
// Output: let person={age:30,city:"New York",name:"John"};console.log(person.name)
```

## Loop Examples

### For Loop

```go
input := `
for (let i = 0; i < 10; i++) {
	console.log(i)
}
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
// Output: for (let i=0;(i<10);i++){console.log(i)}
```

### While Loop

```go
input := `
let count = 0
while (count < 5) {
	console.log("Count:", count)
	count++
}
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
// Output: let count=0;while ((count<5)){console.log("Count:",count);count++}
```

## Function Expressions

```go
input := `
let multiply = function(a, b) {
	return a * b
}
let result = multiply(5, 3)
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
// Output: let multiply=function(a,b){return (a*b)};let result=multiply(5,3)
```

## Error Handling Example

```go
package main

import (
	"fmt"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func parseWithErrorHandling(input string) {
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}
	
	// Check for parser errors
	if errors := p.Errors(); len(errors) > 0 {
		fmt.Println("Parser errors found:")
		for _, e := range errors {
			fmt.Printf("  Line %d, Column %d: %s\n", 
				e.Position.Line, e.Position.Column, e.Message)
		}
		return
	}
	
	fmt.Println("Success:", program.String())
}

func main() {
	// Valid input
	parseWithErrorHandling(`let x = 5`)
	// Output: Success: let x=5
	
	// Invalid input
	parseWithErrorHandling(`let x = 5 +`)
	// Output: Parse error: parse error at line 1, column 11: unexpected token EOF
}
```

## Working with Different Data Types

```go
input := `
// Numbers
let integer = 42
let float = 3.14

// Strings
let single = 'Hello'
let double = "World"

// Booleans
let isTrue = true
let isFalse = false

// Null and undefined
let nothing = null
let notDefined = undefined

// Expressions
let calculation = (10 + 5) * 2
let comparison = age > 18 && hasLicense
`

l := lexer.New(input)
p := parser.New(l)
program, _ := p.ParseProgram()

fmt.Println(program.String())
```

## Practical Use Cases

### 1. Code Formatter

```go
package main

import (
	"fmt"
	"os"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func formatJavaScript(filename string) error {
	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	// Parse and reformat
	l := lexer.New(string(content))
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return err
	}
	
	// Write formatted output
	formatted := program.String()
	return os.WriteFile(filename+".formatted", []byte(formatted), 0644)
}
```

### 2. Syntax Validator

```go
func validateSyntax(code string) (bool, []string) {
	l := lexer.New(code)
	p := parser.New(l)
	
	_, err := p.ParseProgram()
	if err != nil {
		return false, []string{err.Error()}
	}
	
	var errorMessages []string
	for _, e := range p.Errors() {
		errorMessages = append(errorMessages, e.Message)
	}
	
	return len(errorMessages) == 0, errorMessages
}

func main() {
	code := `let x = 5; console.log(x)`
	
	isValid, errors := validateSyntax(code)
	if isValid {
		fmt.Println("Code is valid!")
	} else {
		fmt.Println("Syntax errors found:")
		for _, err := range errors {
			fmt.Println("  -", err)
		}
	}
}
```

### 3. AST Inspector

```go
package main

import (
	"fmt"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func inspectAST(input string) {
	l := lexer.New(input)
	p := parser.New(l)
	program, _ := p.ParseProgram()
	
	fmt.Println("AST Analysis:")
	fmt.Printf("Total statements: %d\n", len(program.Statements))
	
	for i, stmt := range program.Statements {
		fmt.Printf("Statement %d: ", i+1)
		switch s := stmt.(type) {
		case *ast.LetStatement:
			fmt.Printf("Let declaration - Variable: %s\n", s.Name.Value)
		case *ast.FunctionDeclaration:
			fmt.Printf("Function declaration - Name: %s, Params: %d\n", 
				s.Name.Value, len(s.Parameters))
		case *ast.ExpressionStatement:
			fmt.Println("Expression statement")
		case *ast.ReturnStatement:
			fmt.Println("Return statement")
		default:
			fmt.Printf("Other: %T\n", s)
		}
	}
}

func main() {
	code := `
		let name = "XJS"
		function greet() {
			return "Hello!"
		}
		greet()
	`
	
	inspectAST(code)
}
```

## Performance Tips

1. **Reuse Parser Instance**: For multiple parses, create parser once and reset with new lexer
2. **Batch Processing**: Process multiple files in a single run
3. **Error Early**: Check for parsing errors before processing AST

```go
// Efficient batch processing
func parseMultipleFiles(filenames []string) {
	for _, filename := range filenames {
		content, _ := os.ReadFile(filename)
		l := lexer.New(string(content))
		p := parser.New(l)
		
		program, err := p.ParseProgram()
		if err != nil {
			fmt.Printf("Error in %s: %v\n", filename, err)
			continue
		}
		
		// Process the program...
		fmt.Printf("Processed %s successfully\n", filename)
	}
}
```

---

Next: Learn how to extend XJS with [Middleware System](middleware-system.md)