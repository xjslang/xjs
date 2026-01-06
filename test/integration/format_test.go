package integration

import (
	"fmt"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func Example_inlineComments() {
	input := `
	console.log('Hello, World!') // prints a message
	console.log('Bye, bye'); // prints another message`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, _ := p.ParseProgram()
	c := compiler.New().WithPrettyPrint().Compile(program)
	fmt.Println(c.Code)
	// Output:
	// console.log("Hello, World!"); // prints a message
	// console.log("Bye, bye"); // prints another message
}

func Example_format() {
	input := `
		// init vars
		let x = 5
		let y = 10.5
		let name = "Hello World"

		// push items
		let items = []
		items.push(function () {
			console.log("new item")
		})

		// add function
		function add(a, b) {
			return a + b
		}

		// conditional
		if (x < y) {
			console.log("x is less than y")
		}

		// init vars
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
	result := compiler.New().WithPrettyPrint().Compile(program)
	fmt.Println(result.Code)
	// Output:
	// // init vars
	// let x = 5;
	// let y = 10.5;
	// let name = "Hello World";
	//
	// // push items
	// let items = [];
	// items.push(function() {
	//   console.log("new item");
	// }
	// );
	//
	// // add function
	// function add(a, b) {
	//   return a + b;
	// }
	//
	// // conditional
	// if (x < y) {
	//   console.log("x is less than y");
	// }
	//
	// // init vars
	// let numbers = [1, 2, 3, 4, 5];
	// let person = {age: 30, name: "John"};
}
