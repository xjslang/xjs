package parser_test

import (
	"fmt"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/printer"
)

func Example_basic() {
	result, err := testutil.Parse(`function hello() {
	let x = 100
	let y = 200
}`)
	if err != nil {
		panic(err)
	}

	pr := printer.Printer{}
	pr.Init()
	pr.Print(result)
	fmt.Print(pr.String())
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}
