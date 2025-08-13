// Package main demonstrates basic usage of the xjslang library
package main

import (
	"fmt"

	"github.com/xjslang/xjslang/lexer"
	"github.com/xjslang/xjslang/parser"
	"github.com/xjslang/xjslang/token"
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
		
		let numbers = [1, 2, 3, 4, 5]
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

	fmt.Println("\n=== PARSER OUTPUT ===")
	l2 := lexer.New(input)
	p := parser.New(l2)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range p.Errors() {
			fmt.Println("\t" + err)
		}
		return
	}

	fmt.Println("AST:")
	fmt.Println(program.String())
}
