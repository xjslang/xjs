// Package main demonstrates basic usage of the xjs library with function expressions
package main

import (
	"fmt"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func Example_general() {
	input := `
		let x = 5
		let y = 10.5
		let name = "Hello World"

		let items = []
		items.push(function () {
			console.log("new item")
		})

		function add(a, b) {
			return a + b
		}

		if (x < y) {
			console.log("x is less than y")
		}

		let numbers = [1, 2, 3, 4, 5]
		let person = {name: "John", age: 30}
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range p.Errors() {
			fmt.Println("\t" + err)
		}
		return
	}

	fmt.Println(program.String())
	// Output: let x = 5let y = 10.5let name = "Hello World"let items = []items.push(function() {console.log("new item")})function add(a, b) {return (a + b)}if ((x < y)) {console.log("x is less than y")}let numbers = [1, 2, 3, 4, 5]let person = {name: "John", age: 30}
}
