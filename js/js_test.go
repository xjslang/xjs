package js_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/source"
	"github.com/xjslang/xjs/token"
)

func nodeString(node ast.Node) string {
	indentLevel := 0
	var print func(node ast.Node) string
	print = func(node ast.Node) string {
		s := &strings.Builder{}
		indentLevel++
		defer func() {
			indentLevel--
		}()
		indent := strings.Repeat("\t", indentLevel)
		fmt.Fprintf(s, "%T", node)
		switch v := node.(type) {
		case *js.BlockStatement:
			for _, stmt := range v.Statements {
				fmt.Fprintf(s, "\n%s%s", indent, print(stmt))
			}
		case *js.LetStatement:
			fmt.Fprintf(s, "\n%sName: %s", indent, v.Name.Literal)
			fmt.Fprintf(s, "\n%sValue: %s", indent, print(v.Value))
		case *js.FunctionDeclaration:
			fmt.Fprintf(s, "\n%sName: %s", indent, v.Name.Literal)
			fmt.Fprintf(s, "\n%sBody: %s", indent, print(v.Body))
		case *ast.GroupedExpression:
			fmt.Fprintf(s, "\n%sValue: %s", indent, print(v.Value))
		case *js.InfixOperator:
			fmt.Fprintf(s, "\n%sLeftValue: %s", indent, print(v.LeftValue))
			fmt.Fprintf(s, "\n%sOperator: %q", indent, v.Operator.Type.String())
			fmt.Fprintf(s, "\n%sRightValue: %s", indent, print(v.RightValue))
		case *ast.IntegerLiteral:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		case *ast.StringLiteral:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		case *ast.BooleanLiteral:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		case *ast.Identifier:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		}
		return s.String()
	}
	return print(node)
}

func ExampleParse() {
	result, err := js.Parse([]byte(`function hello() {
	let x = 100
	let y = 200
}`))
	if err != nil {
		panic(err)
	}

	pr := printer.New()
	result.PrintTo(pr)
	fmt.Print(pr.String())
	// Output:
	// function hello() {
	//   let x = 100;
	//   let y = 200;
	// }
}

func TestParser(t *testing.T) {
	input := `
		function printHello() {
			let x = 100;
			let y = 200;
		}

		let x = 100;
		let y = 200;`
	pr, err := js.Parse([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	expected := `*js.BlockStatement
	*js.FunctionDeclaration
		Name: printHello
		Body: *js.BlockStatement
			*js.LetStatement
				Name: x
				Value: *ast.IntegerLiteral{Value: "100"}
			*js.LetStatement
				Name: y
				Value: *ast.IntegerLiteral{Value: "200"}
	*js.LetStatement
		Name: x
		Value: *ast.IntegerLiteral{Value: "100"}
	*js.LetStatement
		Name: y
		Value: *ast.IntegerLiteral{Value: "200"}`
	if result := nodeString(pr); result != expected {
		t.Errorf("Invalid node:\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestKeepParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedAST    string
		expectedErrors int
		expectedMsg    string
	}{
		{
			input: `
			let x = 100
			aaa // unknown statement
			let y = 200
			bbb // unknown statement
			let z = 300`,
			expectedAST: `*js.BlockStatement
	*js.LetStatement
		Name: x
		Value: *ast.IntegerLiteral{Value: "100"}
	*js.LetStatement
		Name: y
		Value: *ast.IntegerLiteral{Value: "200"}
	*js.LetStatement
		Name: z
		Value: *ast.IntegerLiteral{Value: "300"}`,
			expectedErrors: 2,
			expectedMsg:    "Unknown statement",
		},
		{
			input: `
			let x = 100
			/* unknown statement at the end of file */ aaa`,
			expectedAST: `*js.BlockStatement
	*js.LetStatement
		Name: x
		Value: *ast.IntegerLiteral{Value: "100"}`,
			expectedErrors: 1,
			expectedMsg:    "Unknown statement",
		},
		{
			input: "let x = 100; /* unknown statement at the end of line */ aaa",
			expectedAST: `*js.BlockStatement
	*js.LetStatement
		Name: x
		Value: *ast.IntegerLiteral{Value: "100"}`,
			expectedErrors: 1,
			expectedMsg:    "Unknown statement",
		},
		{
			input:          "let x = ;",
			expectedAST:    "*js.BlockStatement",
			expectedErrors: 1,
			expectedMsg:    "Expected value",
		},
		{
			input: `
			let x = 100; let y =
			let z = 200`,
			expectedAST: `*js.BlockStatement
	*js.LetStatement
		Name: x
		Value: *ast.IntegerLiteral{Value: "100"}
	*js.LetStatement
		Name: z
		Value: *ast.IntegerLiteral{Value: "200"}`,
			expectedErrors: 1,
			expectedMsg:    "Expected value",
		},
	}
	for i, test := range tests {
		for j := range 2 {
			testname := fmt.Sprintf("test %d", i)
			if j > 0 {
				testname += " function wrap"
			}
			t.Run(testname, func(t *testing.T) {
				input := test.input
				if j > 0 {
					input = fmt.Sprintf("function main(){\n%s}", input)
				}
				pr, err := js.Parse([]byte(input))
				if j == 0 {
					if result := nodeString(pr); result != test.expectedAST {
						t.Fatalf("Invalid node:\nExpected:\n%s\nGot:\n%s", test.expectedAST, result)
					}
				}
				list, ok := err.(parser.ErrorList)
				if !ok {
					t.Fatalf("Expected %T, got %T", list, err)
				}
				if n := len(list); n != test.expectedErrors {
					t.Fatalf("Expected %d errors, got %d", test.expectedErrors, n)
				}
				for _, err := range list {
					if msg := err.Message; msg != test.expectedMsg {
						t.Fatalf("Expected %q, got %q", test.expectedMsg, msg)
					}
				}
			})
		}
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected parser.Error
	}{
		{"missing token", "function\nhello({}", parser.Error{
			Range: source.Range{
				Start: source.Position{Line: 1, Column: 6},
				End:   source.Position{Line: 1, Column: 7},
			},
			Message: "Expected )"},
		},
		{"missing semicolon", "let x = 100 let y = 200", parser.Error{
			Range: source.Range{
				Start: source.Position{Line: 0, Column: 12},
				End:   source.Position{Line: 0, Column: 15},
			},
			Message: "Expected statement terminator"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := js.Parse([]byte(test.input))
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			list, ok := err.(parser.ErrorList)
			if !ok {
				t.Fatalf("Expected ErrorList, got %T", err)
			}
			if n := len(list); n != 1 {
				t.Errorf("Expected one error, got %d", n)
			} else if list[0].Message != test.expected.Message {
				t.Errorf("Expected error %q, got %q", test.expected.Message, list[0].Message)
			} else if list[0].Range != test.expected.Range {
				t.Error("position error")
			}
		})
	}

	t.Run("multiple errors", func(t *testing.T) {
		expectedErrors := []string{
			"Expected value",
			"Expected =",
			"Expected statement terminator",
			"Expected )",
		}
		_, err := js.Parse([]byte(`
		let x = !
		let y
		let x = 100.

		function print(`))
		errors, ok := err.(parser.ErrorList)
		if !ok {
			t.Fatalf("Expected ErrorList, got %T", err)
		}
		if n := len(errors); n != len(expectedErrors) {
			t.Fatalf("Expected %d errors, got %d", len(expectedErrors), n)
		}
		for i, expectedError := range expectedErrors {
			if errors[i].Message != expectedError {
				t.Errorf("Expected %q, got %q", expectedError, errors[i].Message)
			}
		}
	})
}

type VarStatement struct {
	Name  token.Token
	Value ast.Expression
}

func (ls *VarStatement) PrintTo(p *printer.Printer) {
	p.PrintString("var ")
	p.PrintString(ls.Name.Literal)
	p.PrintString(" = ")
	ls.Value.PrintTo(p)
	p.PrintRune(';')
}

func TestUseStatementParser(t *testing.T) {
	l := &lexer.Lexer{}
	l.Init([]byte("var x = 100"))
	b := js.Builder{}
	b.InstallCorePlugins()
	b.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "var" {
			p.AdvanceToken() // consume "var"
			name := p.CurrentToken
			if err := p.Expect(token.IDENT); err != nil {
				return nil, err
			}
			if err := p.Expect(token.ASSIGN); err != nil {
				return nil, err
			}
			value, err := p.ParseExpression()
			if err != nil {
				return nil, err
			}
			if err := js.ExpectSemi(p); err != nil {
				return nil, err
			}
			return &VarStatement{Name: name, Value: value}, nil
		}
		return next()
	})
	p := b.Build(l)

	prog, err := js.ParseProgram(p)
	if err != nil {
		t.Fatal(err)
	}
	prt := printer.New()
	prog.PrintTo(prt)
	expected := "var x = 100;\n"
	if result := prt.String(); result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

type DeferStatement struct {
	//
}

func (node *DeferStatement) PrintTo(p *printer.Printer) {
	//
}

func DeferPlugin(b *parser.Builder) {
	b.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "defer" {
			p.AdvanceToken() // consume defer
			if p.CurrentToken.Type != token.IDENT {
				msg := "Expected " + token.IDENT.String()
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			p.AdvanceToken()
			if p.CurrentToken.Type != token.LPAREN {
				msg := "Expected " + token.LPAREN.String()
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			p.AdvanceToken()
			if p.CurrentToken.Type != token.RPAREN {
				msg := "Expected " + token.RPAREN.String()
				p.AddError(msg)
				return nil, errors.New(msg)
			}
			p.AdvanceToken()
			return &DeferStatement{}, nil
		}
		return next()
	})
}

func TestBuilder(t *testing.T) {
	l := &lexer.Lexer{}
	l.Init([]byte(`defer dbClose()`))

	b := js.Builder{}
	b.InstallCorePlugins()
	b.Install(DeferPlugin)
	p := b.Build(l)
	pr, err := js.ParseProgram(p)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(pr.Statements); n != 1 {
		t.Fatalf("Expected 1 statements, got %d", n)
	}
	stmt := pr.Statements[0]
	if _, ok := stmt.(*DeferStatement); !ok {
		t.Fatalf("Expected *DeferStatement, got %T", stmt)
	}
}

type PowExpression struct {
	LeftValue  ast.Expression
	Operator   token.Token
	RightValue ast.Expression
}

func (node *PowExpression) PrintTo(p *printer.Printer) {
	p.PrintString("Math.pow(")
	node.LeftValue.PrintTo(p)
	p.PrintString(", ")
	node.RightValue.PrintTo(p)
	p.PrintRune(')')
}

func TestRegisterInfixOperator(t *testing.T) {
	l := &lexer.Lexer{}
	powType := token.RegisterType("**")
	l.Init([]byte("1 + 2 ** 3"))

	b := js.Builder{}
	b.InstallCorePlugins()
	b.UseTokenizer(func(l *lexer.Lexer, next func() token.Token) token.Token {
		if l.CurrentChar == '*' && l.PeekChar() == '*' {
			l.AdvanceChar()
			l.AdvanceChar()
			return token.Token{Type: powType, Literal: powType.String()}
		}
		return next()
	})
	b.RegisterInfixOperator(powType, 3, func(op token.Token, left, right ast.Expression) ast.Expression {
		return &PowExpression{LeftValue: left, Operator: op, RightValue: right}
	})
	p := b.Build(l)
	exp, err := p.ParseExpression()
	if err != nil {
		t.Fatal(err)
	}

	expected := `*js.InfixOperator
	LeftValue: *ast.IntegerLiteral{Value: "1"}
	Operator: "+"
	RightValue: *js_test.PowExpression`
	if result := nodeString(exp); result != expected {
		t.Errorf("Invalid node:\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single value", "3", `*ast.IntegerLiteral{Value: "3"}`},
		{
			name:  "boolean values",
			input: "true + false",
			expected: `*js.InfixOperator
	LeftValue: *ast.BooleanLiteral{Value: "true"}
	Operator: "+"
	RightValue: *ast.BooleanLiteral{Value: "false"}`,
		},
		{
			name:  "identifiers",
			input: "x + 2*y",
			expected: `*js.InfixOperator
	LeftValue: *ast.Identifier{Value: "x"}
	Operator: "+"
	RightValue: *js.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "2"}
		Operator: "*"
		RightValue: *ast.Identifier{Value: "y"}`,
		},
		{
			name:  "basic",
			input: "1 + 2 - 3",
			expected: `*js.InfixOperator
	LeftValue: *js.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "1"}
		Operator: "+"
		RightValue: *ast.IntegerLiteral{Value: "2"}
	Operator: "-"
	RightValue: *ast.IntegerLiteral{Value: "3"}`,
		},
		{
			name:  "complex",
			input: "5 - 2 * '3' % 5 + 1",
			expected: `*js.InfixOperator
	LeftValue: *js.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "5"}
		Operator: "-"
		RightValue: *js.InfixOperator
			LeftValue: *js.InfixOperator
				LeftValue: *ast.IntegerLiteral{Value: "2"}
				Operator: "*"
				RightValue: *ast.StringLiteral{Value: "'3'"}
			Operator: "%"
			RightValue: *ast.IntegerLiteral{Value: "5"}
	Operator: "+"
	RightValue: *ast.IntegerLiteral{Value: "1"}`,
		},
		{
			name:  "parentheses",
			input: "2 * (3 + 5) - (1)",
			expected: `*js.InfixOperator
	LeftValue: *js.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "2"}
		Operator: "*"
		RightValue: *ast.GroupedExpression
			Value: *js.InfixOperator
				LeftValue: *ast.IntegerLiteral{Value: "3"}
				Operator: "+"
				RightValue: *ast.IntegerLiteral{Value: "5"}
	Operator: "-"
	RightValue: *ast.GroupedExpression
		Value: *ast.IntegerLiteral{Value: "1"}`,
		},
		{
			name: "with comments",
			input: `2 // first operand
			/ 3 /* second operand */
			* 4`,
			expected: `*js.InfixOperator
	LeftValue: *js.InfixOperator
		LeftValue: *ast.IntegerLiteral{Value: "2"}
		Operator: "/"
		RightValue: *ast.IntegerLiteral{Value: "3"}
	Operator: "*"
	RightValue: *ast.IntegerLiteral{Value: "4"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer.Lexer{}
			l.Init([]byte(test.input))
			b := js.Builder{}
			b.InstallCorePlugins()
			p := b.Build(l)
			exp, err := p.ParseExpression()
			if err != nil {
				t.Fatal(err)
			}
			if result := nodeString(exp); result != test.expected {
				t.Errorf("Invalid node:\nExpected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}

func TestMalformedExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "missing value",
			input:    "1 + 2*",
			expected: "Expected value",
		},
		{
			name:     "missing right parenthesis",
			input:    "2 * (3 + 5",
			expected: "Expected )",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := &lexer.Lexer{}
			l.Init([]byte(test.input))
			b := js.Builder{}
			b.InstallCorePlugins()
			p := b.Build(l)
			_, err := p.ParseExpression()
			if err == nil {
				t.Fatal("An error was expected")
			}
			if result := err.Error(); result != test.expected {
				t.Errorf("Expected error to be %q, got %q", test.expected, result)
			}
		})
	}
}
