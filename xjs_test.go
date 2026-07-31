package xjs_test

import (
	"fmt"

	"github.com/xjslang/xjs"
)

func Example_basic() {
	input := `function hello() {
	let x = 100
	let y = 200
}`
	result, err := xjs.Parse([]byte(input))
	if err != nil {
		panic(err)
	}
	out, err := xjs.Print(result)
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
