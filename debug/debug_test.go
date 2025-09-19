package debug

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func TestDebugToString(t *testing.T) {
	stmt := &ast.LetStatement{
		Name:  &ast.Identifier{Value: "x"},
		Value: &ast.IntegerLiteral{Token: token.Token{Literal: "5"}},
	}

	result := ToString(stmt)
	expected := "let x=5"

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func ExampleToString() {
	stmt := &ast.LetStatement{
		Name:  &ast.Identifier{Value: "x"},
		Value: &ast.IntegerLiteral{Token: token.Token{Literal: "5"}},
	}

	fmt.Println(ToString(stmt))
	// Output: let x=5
}

func ExamplePrint() {
	stmt := &ast.LetStatement{
		Name:  &ast.Identifier{Value: "x"},
		Value: &ast.IntegerLiteral{Token: token.Token{Literal: "5"}},
	}

	Print(stmt)
	// Output:
	// (*ast.LetStatement)({
	//    Token: (token.Token) {
	//       Type: (token.Type) 0,
	//       Literal: (string) "",
	//       Line: (int) 0,
	//       Column: (int) 0
	//    },
	//    Name: (*ast.Identifier)({
	//       Token: (token.Token) {
	//          Type: (token.Type) 0,
	//          Literal: (string) "",
	//          Line: (int) 0,
	//          Column: (int) 0
	//       },
	//       Value: (string) (len=1) "x"
	//    }),
	//    Value: (*ast.IntegerLiteral)({
	//       Token: (token.Token) {
	//          Type: (token.Type) 0,
	//          Literal: (string) (len=1) "5",
	//          Line: (int) 0,
	//          Column: (int) 0
	//       }
	//    })
	// })
}
