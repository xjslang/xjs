package debug

import (
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func TestToString(t *testing.T) {
	tests := []struct{ name, input, expected string }{
		{"LetStatement", "let x = 5", "let x=5"},
		{"FunctionDeclaration", "function add(a, b){ return a+b }", "function add(a,b){return (a+b)}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			stmt := p.ParseStatement()
			if ToString(stmt) != tt.expected {
				t.Errorf("ParseStatement() got %q, want %q", ToString(stmt), tt.expected)
			}
		})
	}
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
