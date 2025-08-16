/*
Package token defines the token types and structures used by the XJS lexer and parser.

This package provides all the token types supported by the XJS language, including
literals, operators, keywords, and delimiters. It also includes utility functions for
working with tokens and identifiers.

# Token Types

The following token types are supported:

- Literals: IDENT, INT, FLOAT, STRING
- Operators: +, -, *, /, %, ==, !=, ===, !==, <, >, <=, >=, &&, ||, !, ++, --
- Keywords: function, let, if, else, while, for, return, true, false, null
- Delimiters: (, ), {, }, [, ], ,, ;, :, .

Example:

	token := token.Token{
		Type: token.IDENT,
		Literal: "myVariable",
		Line: 1,
		Column: 5,
	}

	fmt.Println(token.String())
	// Output: {Type: IDENT, Literal: "myVariable", Line: 1, Col: 5}
*/
package token
