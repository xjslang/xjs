//go:build integration

// Package main demonstrates basic usage of the xjs library with function expressions
package integration

import (
	"fmt"

	"github.com/xjslang/xjs/compiler"
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
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("ParseProgram error = %v\n", err)
		return
	}
	result := compiler.New().Compile(program)
	fmt.Println(result.Code)
	// Output: let x=5;let y=10.5;let name="Hello World";let items=[];items.push(function(){console.log("new item");});function add(a,b){return a+b;}if(x<y){console.log("x is less than y");}let numbers=[1,2,3,4,5];let person={name:"John",age:30};
}
