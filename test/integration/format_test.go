package integration

import (
	"fmt"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func Example_objectLiteral() {
	input := `
	let user = {
		name: 'John Smith', // user name
		age: 35, // user age
	}`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, _ := p.ParseProgram()
	c := compiler.New().WithPrettyPrint().Compile(program)
	fmt.Println(c.Code)
	// Output:
	// let user = {
	//   name: "John Smith", // user name
	//   age: 35 // user age
	// };
}

func Example_inlineComments() {
	input := `
	console.log('Hello, World!') // prints a message
	console.log('Bye, bye'); // prints another message
	let x = 5 // variable without semicolon
	let y = 10; // variable with semicolon
	function test() {
		return x // return without semicolon
	}
	function test2() {
		return y; // return with semicolon
	}`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, _ := p.ParseProgram()
	c := compiler.New().WithPrettyPrint().Compile(program)
	fmt.Println(c.Code)
	// Output:
	// console.log("Hello, World!"); // prints a message
	// console.log("Bye, bye"); // prints another message
	// let x = 5; // variable without semicolon
	// let y = 10; // variable with semicolon
	// function test() {
	//   return x; // return without semicolon
	// }
	// function test2() {
	//   return y; // return with semicolon
	// }
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
	// let person = {name: "John", age: 30};
}
