// Package xjs provides a JavaScript-like language lexer and parser.
// This package offers clean, simplified JavaScript parsing without redundant features
// like classes, arrow functions, const/var declarations, etc.
//
// Example usage:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/xjslang/xjs/lexer"
//		"github.com/xjslang/xjs/parser"
//	)
//
//	func main() {
//		input := `let x = 5; function add(a, b) { return a + b; }`
//
//		l := lexer.New(input)
//		p := parser.New(l)
//		program := p.ParseProgram()
//
//		fmt.Println(program.String())
//	}
package xjs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

// Parse is a convenience function that lexes and parses the given input
// and returns the resulting AST program node along with any parsing errors.
func Parse(input string) (*ast.Program, []string) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return program, p.Errors()
}

// Version returns the current version of xjslang
const Version = "0.1.0"
