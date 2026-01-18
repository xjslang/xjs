package compiler

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestComments(t *testing.T) {
	input := `(
// function expression comment
function(){
})`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram error %v", err)
	}
	result := New().WithPrettyPrint(WithSemi(false)).Compile(program)
	fmt.Println(result.Code)
}
