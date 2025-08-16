/*
Package lexer provides lexical analysis functionality for the XJS language.

The lexer tokenizes source code into a sequence of tokens that can be consumed
by the parser. It supports all XJS language features including variables,
functions, control flow, arrays, objects, and various operators.

Example:

	input := `let x = 5; function add(a, b) { return a + b; }`

	l := lexer.New(input)
	for {
		tok := l.NextToken()
		fmt.Println(tok)
		if tok.Type == token.EOF {
			break
		}
	}

# Position Tracking

The lexer automatically tracks line and column positions for each token,
which is useful for error reporting and debugging.

# Supported Features

- Variable declarations (let)
- Function declarations
- All arithmetic and logical operators
- String literals with single or double quotes
- Integer and floating-point numbers
- Keywords and identifiers
- Arrays and objects
- Control flow statements
*/
package lexer
