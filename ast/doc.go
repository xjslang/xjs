/*
Package ast defines the Abstract Syntax Tree nodes for the XJS language.

This package provides interfaces and concrete types representing different language
constructs such as statements, expressions, and literals. The AST is designed to
be simple and easy to work with while covering all essential JavaScript-like
language features.

Node Types

The AST defines three main categories of nodes:

- Statements: let, function declarations, if/else, while, for, return, blocks
- Expressions: literals, binary/unary operations, calls, member access
- Literals: integers, floats, strings, booleans, null, arrays, objects

Example:

	// Creating an AST node for: let x = 5
	letStmt := &ast.LetStatement{
		Token: token.Token{Type: token.LET, Literal: "let"},
		Name: &ast.Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "x"},
			Value: "x",
		},
		Value: &ast.IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "5"},
			Value: 5,
		},
	}

	fmt.Println(letStmt.String()) // Output: let x = 5;

String Representation

All AST nodes implement the String() method, providing a readable representation
of the parsed code structure.
*/
package ast
