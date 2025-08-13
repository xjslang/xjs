/*
Package parser provides syntax analysis functionality for the xjslang language.

The parser uses Pratt parsing (also known as "top-down operator precedence parsing")
to build an Abstract Syntax Tree (AST) from tokens provided by the lexer. It handles
operator precedence correctly and provides detailed error reporting with line and
column information.

Example:

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

	fmt.Println(program.String())

Error Handling

The parser provides comprehensive error reporting with specific line and column
information for each error encountered during parsing.

Operator Precedence

The parser correctly handles operator precedence for all supported operators:

1. Assignment (=)
2. Logical OR (||)
3. Logical AND (&&)
4. Equality (==, !=, ===, !==)
5. Comparison (<, >, <=, >=)
6. Addition/Subtraction (+, -)
7. Multiplication/Division/Modulo (*, /, %)
8. Unary operators (!, -, ++, --)
9. Function calls and member access

Supported Language Features

- Variable declarations with let
- Function declarations
- Control flow: if/else, while, for loops
- Expressions: arithmetic, logical, comparison
- Literals: numbers, strings, booleans, null
- Data structures: arrays, objects
- Member access: dot notation and computed access
- Function calls with arguments
*/
package parser
