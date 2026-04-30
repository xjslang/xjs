package ast

import (
	"testing"

	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

func TestPrintTo(t *testing.T) {
	program := &BlockStatement{
		Statements: []Statement{
			&FunctionDeclaration{
				Name: token.Token{Type: token.IDENT, Literal: "foo"},
				Body: &BlockStatement{
					Statements: []Statement{
						&FunctionDeclaration{
							Name: token.Token{Type: token.IDENT, Literal: "boo"},
							Body: &BlockStatement{
								Statements: []Statement{
									&LetStatement{Name: token.Token{Type: token.IDENT, Literal: "m"}, Value: &StringLiteral{Value: "'mmm'"}},
									&LetStatement{Name: token.Token{Type: token.IDENT, Literal: "n"}, Value: &StringLiteral{Value: "'nnn'"}},
								},
							},
						},
						&LetStatement{
							Name:  token.Token{Type: token.IDENT, Literal: "a"},
							Value: &StringLiteral{Value: "'aaa'"},
						},
						&LetStatement{
							Name:  token.Token{Type: token.IDENT, Literal: "b"},
							Value: &StringLiteral{Value: "'bbb'"},
						},
					},
				},
			},
			&LetStatement{
				Name:  token.Token{Type: token.IDENT, Literal: "x"},
				Value: &IntegerLiteral{Value: "100"},
			},
			&LetStatement{
				Name:  token.Token{Type: token.IDENT, Literal: "y"},
				Value: &IntegerLiteral{Value: "200"},
			},
		},
	}
	pr := printer.New(printer.WithIndent("\t"))
	program.PrintTo(pr)
	expected := "function foo() {\n\tfunction boo() {\n\t\tlet m = 'mmm';\n\t\tlet n = 'nnn';\n\t}\n\tlet a = 'aaa';\n\tlet b = 'bbb';\n}\nlet x = 100;\nlet y = 200;\n"
	if result := pr.String(); result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
