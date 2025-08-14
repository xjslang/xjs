// Package main demonstrates the xjs library functionality
package main

import (
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func main() {
	input := `
		let x = 5
		let y = 10.5
		let name = "Hello World"
		
		function add(a, b) {
			return a + b
		}
		
		if (x < y) {
			console.log("x is less than y")
		}
		
		let numbers = [1, 2, 3]
		let person = {name: "John", age: 30}
	`

	fmt.Println("=== LEXER OUTPUT ===")
	l := lexer.New(input)

	for {
		tok := l.NextToken()
		fmt.Println(tok)
		if tok.Type == token.EOF {
			break
		}
	}

	fmt.Println("\n=== PARSER OUTPUT (using convenience function) ===")
	program, errors := xjs.Parse(input)

	if len(errors) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range errors {
			fmt.Println("\t" + err)
		}
		return
	}

	fmt.Println("AST:")
	fmt.Println(program.String())

	fmt.Printf("\nxjs version: %s\n", xjs.Version)
}
