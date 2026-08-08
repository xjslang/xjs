package xjs_test

import (
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/js"
)

func Example_basic() {
	input := `function hello() {
		let x = 100
		let y = 200
		}`

	// code --> AST
	sc := xjs.ScannerBuilder().Build([]byte(input))
	p := xjs.ParserBuilder().Build(sc)
	result, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}

	// AST --> code
	pr := xjs.PrinterBuilder().Build()
	pr.Print(result)
	out, err := pr.Output()
	if err != nil {
		panic(err)
	}
	fmt.Print(out)

	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}
