package ast

import (
	"testing"

	"github.com/xjslang/xjs/printer"
)

func TestPrintTo(t *testing.T) {
	program := &BlockStatement{
		Statements: []Statement{
			&FunctionDeclaration{
				Name: &Identifier{Value: "foo"},
				Body: &BlockStatement{
					Statements: []Statement{
						&FunctionDeclaration{
							Name: &Identifier{Value: "boo"},
							Body: &BlockStatement{
								Statements: []Statement{
									&LetStatement{Name: &Identifier{Value: "m"}, Value: &StringLiteral{Value: "'mmm'"}},
									&LetStatement{Name: &Identifier{Value: "n"}, Value: &StringLiteral{Value: "'nnn'"}},
								},
							},
						},
						&LetStatement{
							Name:  &Identifier{Value: "a"},
							Value: &StringLiteral{Value: "'aaa'"},
						},
						&LetStatement{
							Name:  &Identifier{Value: "y"},
							Value: &StringLiteral{Value: "'bbb'"},
						},
					},
				},
			},
			&LetStatement{
				Name:  &Identifier{Value: "x"},
				Value: &IntegerLiteral{Value: "100"},
			},
			&LetStatement{
				Name:  &Identifier{Value: "y"},
				Value: &IntegerLiteral{Value: "200"},
			},
		},
	}
	pr := printer.New(printer.WithIndent("\t"))
	program.PrintTo(pr)
	expected := "function foo() {\n\tfunction boo() {\n\t\tlet m = 'mmm';\n\t\tlet n = 'nnn';\n\t}\n\tlet a = 'aaa';\n\tlet y = 'bbb';\n}\nlet x = 100;\nlet y = 200;\n"
	if result := pr.String(); result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
