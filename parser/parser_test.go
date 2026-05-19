package parser_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
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
	pr.PrintNode(result)
	fmt.Print(pr.String())
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}

func TestExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: "1 + 2 * (3 + 5) - 4",
			expected: `BinaryExpr
	LeftValue: BinaryExpr
		LeftValue: Integer{Value: "1"}
		Operator: "+"
		RightValue: BinaryExpr
			LeftValue: Integer{Value: "2"}
			Operator: "*"
			RightValue: ParenExpr
				Value: BinaryExpr
					LeftValue: Integer{Value: "3"}
					Operator: "+"
					RightValue: Integer{Value: "5"}
	Operator: "-"
	RightValue: Integer{Value: "4"}`,
		},
		{
			input: "foo() + 1",
			expected: `BinaryExpr
	LeftValue: CallExpr
		Function: Ident{Value: "foo"}
	Operator: "+"
	RightValue: Integer{Value: "1"}`,
		},
		{
			input: "foo(1, 2, 3)",
			expected: `CallExpr
	Function: Ident{Value: "foo"}
	Arguments[0]: Integer{Value: "1"}
	Arguments[1]: Integer{Value: "2"}
	Arguments[2]: Integer{Value: "3"}`,
		},
		{
			input: "2 * (pow(2, 1 + 3) + 4)",
			expected: `BinaryExpr
	LeftValue: Integer{Value: "2"}
	Operator: "*"
	RightValue: ParenExpr
		Value: BinaryExpr
			LeftValue: CallExpr
				Function: Ident{Value: "pow"}
				Arguments[0]: Integer{Value: "2"}
				Arguments[1]: BinaryExpr
					LeftValue: Integer{Value: "1"}
					Operator: "+"
					RightValue: Integer{Value: "3"}
			Operator: "+"
			RightValue: Integer{Value: "4"}`,
		},
	}
	for i, test := range tests {
		t.Run("exp "+strconv.Itoa(i), func(t *testing.T) {
			sc := &scanner.Scanner{}
			sc.Init([]byte(test.input))
			p := &parser.Parser{}
			p.Init(sc)
			result, err := p.ParseExpression()
			if err != nil {
				t.Fatal(err)
			}
			if got := testutil.NodeString(result); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}
