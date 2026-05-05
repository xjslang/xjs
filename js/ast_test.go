package js

import (
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

func TestPrintTo(t *testing.T) {
	program := &BlockStatement{
		Statements: []ast.Statement{
			&FunctionDeclaration{
				Name: token.Token{Type: token.IDENT, Literal: "foo"},
				Body: &BlockStatement{
					Statements: []ast.Statement{
						&FunctionDeclaration{
							Name: token.Token{Type: token.IDENT, Literal: "boo"},
							Body: &BlockStatement{
								Statements: []ast.Statement{
									&LetStatement{Name: token.Token{Type: token.IDENT, Literal: "m"}, Value: &ast.StringLiteral{Value: "'mmm'"}},
									&LetStatement{Name: token.Token{Type: token.IDENT, Literal: "n"}, Value: &ast.StringLiteral{Value: "'nnn'"}},
								},
							},
						},
						&LetStatement{
							Name:  token.Token{Type: token.IDENT, Literal: "a"},
							Value: &ast.BooleanLiteral{Value: "true"},
						},
						&LetStatement{
							Name:  token.Token{Type: token.IDENT, Literal: "b"},
							Value: &ast.BooleanLiteral{Value: "false"},
						},
						&LetStatement{
							Name:  token.Token{Type: token.IDENT, Literal: "c"},
							Value: &ast.Identifier{Value: "x"},
						},
					},
				},
			},
			&LetStatement{
				Name:  token.Token{Type: token.IDENT, Literal: "x"},
				Value: &ast.IntegerLiteral{Value: "100"},
			},
			&LetStatement{
				Name:  token.Token{Type: token.IDENT, Literal: "y"},
				Value: &ast.IntegerLiteral{Value: "200"},
			},
		},
	}
	pr := printer.New(printer.WithIndent("\t"))
	program.PrintTo(pr)
	expected := "function foo() {\n\tfunction boo() {\n\t\tlet m = 'mmm';\n\t\tlet n = 'nnn';\n\t}\n\tlet a = true;\n\tlet b = false;\n\tlet c = x;\n}\nlet x = 100;\nlet y = 200;\n"
	if result := pr.String(); result != expected {
		t.Errorf("Invalid node:\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}
