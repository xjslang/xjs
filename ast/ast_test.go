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
							Value: &IntegerLiteral{Value: "'aaa'"},
						},
						&LetStatement{
							Name:  &Identifier{Value: "y"},
							Value: &IntegerLiteral{Value: "'bbb'"},
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

	t.Run("without format", func(t *testing.T) {
		pr := printer.New()
		program.PrintTo(pr)
		expected := "function foo(){function boo(){let m='mmm';let n='nnn';}let a='aaa';let y='bbb';}let x=100;let y=200;"
		if result := pr.String(); result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("with format", func(t *testing.T) {
		pr := printer.New().WithFormat(printer.WithIndent("\t"), printer.WithSemicolon())
		program.PrintTo(pr)
		expected := "function foo() {\n\tfunction boo() {\n\t\tlet m = 'mmm';\n\t\tlet n = 'nnn';\n\t}\n\tlet a = 'aaa';\n\tlet y = 'bbb';\n}\nlet x = 100;\nlet y = 200;\n"
		if result := pr.String(); result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}
