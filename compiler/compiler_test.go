package compiler

import (
	"fmt"
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestComments(m *testing.T) {
	input := `// expression comment
1 + 2`
	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, _ := p.ParseProgram()
	result := New().WithPrettyPrint(WithSemi(false)).Compile(program)
	fmt.Println(result.Code)
}
