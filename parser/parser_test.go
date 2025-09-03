package parser

import (
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
)

func TestParseBasicLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"integer literal", "42", "42"},
		{"float literal", "3.14", "3.14"},
		{"string literal", `"hello"`, `"hello"`},
		{"boolean true", "true", "true"},
		{"boolean false", "false", "false"},
		{"null literal", "null", "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseBinaryExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"addition", "1 + 2", "(1 + 2)"},
		{"subtraction", "5 - 3", "(5 - 3)"},
		{"multiplication", "3 * 4", "(3 * 4)"},
		{"division", "8 / 2", "(8 / 2)"},
		{"modulo", "10 % 3", "(10 % 3)"},
		{"equality", "x == y", "(x === y)"},
		{"inequality", "x != y", "(x !== y)"},
		{"less than", "x < y", "(x < y)"},
		{"greater than", "x > y", "(x > y)"},
		{"logical and", "a && b", "(a && b)"},
		{"logical or", "a || b", "(a || b)"},
		{"operator precedence", "1 + 2 * 3", "(1 + (2 * 3))"},
		{"parentheses override", "(1 + 2) * 3", "((1 + 2) * 3)"},
		{"complex precedence", "1 + 2 * 3 - 4", "((1 + (2 * 3)) - 4)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseUnaryExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"negation", "-x", "(-x)"},
		{"logical not", "!true", "(!true)"},
		{"prefix increment", "++x", "(++x)"},
		{"prefix decrement", "--x", "(--x)"},
		{"postfix increment", "x++", "(x++)"},
		{"postfix decrement", "x--", "(x--)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple identifier", "x", "x"},
		{"camelCase identifier", "myVariable", "myVariable"},
		{"underscore identifier", "_private", "_private"},
		{"dollar sign identifier", "$special", "$special"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseLetStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"let without assignment", "let x", "let x"},
		{"let with assignment", "let x = 5", "let x=5"},
		{"let with expression", "let result = 1 + 2", "let result=(1 + 2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseFunctionDeclarations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"function no params", "function test() {}", "function test(){}"},
		{"function with params", "function add(a, b) {}", "function add(a,b){}"},
		{"function with body", "function greet() { return 'hello' }", "function greet(){return \"hello\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseArrayLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty array", "[]", "[]"},
		{"array with numbers", "[1, 2, 3]", "[1,2,3]"},
		{"array with mixed types", "[1, 'hello', true]", "[1,\"hello\",true]"},
		{"nested arrays", "[[1, 2], [3, 4]]", "[[1,2],[3,4]]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseAssignmentExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple assignment", "x = 5", "x=5"},
		{"compound addition", "x += 10", "x +=10"},
		{"compound subtraction", "x -= 5", "x -=5"},
		{"assignment with expression", "x = 1 + 2", "x=(1 + 2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseCallExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"function call no args", "foo()", "foo()"},
		{"function call with args", "add(1, 2)", "add(1,2)"},
		{"chained calls", "obj.method()", "obj.method()"},
		{"nested calls", "outer(inner())", "outer(inner())"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestParseMemberExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"dot notation", "obj.prop", "obj.prop"},
		{"bracket notation", "obj[key]", "obj[key]"},
		{"chained access", "obj.prop.value", "obj.prop.value"},
		{"mixed access", "obj.prop[0]", "obj.prop[0]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				t.Errorf("Parse(%q) errors = %v", tt.input, p.Errors())
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			if program.String() != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, program.String(), tt.expected)
			}
		})
	}
}

func TestUseStatementHandler(t *testing.T) {
	input := "custom_statement"

	l := lexer.New(input)
	p := New(l)

	handler := func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Literal == "custom_statement" {
			return &ast.ExpressionStatement{
				Expression: &ast.Identifier{
					Token: p.CurrentToken,
					Value: "handled_statement",
				},
			}
		}
		return next()
	}

	p.UseStatementHandler(handler)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("Parse with custom handler errors = %v", p.Errors())
		return
	}

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		return
	}

	if program.String() != "handled_statement" {
		t.Errorf("Expected custom handler result, got %v", program.String())
	}
}

func TestUseExpressionHandler(t *testing.T) {
	input := "special_expr"

	l := lexer.New(input)
	p := New(l)

	handler := func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Literal == "special_expr" {
			return &ast.Identifier{
				Token: p.CurrentToken,
				Value: "handled_expression",
			}
		}
		return next()
	}

	p.UseExpressionHandler(handler)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("Parse with custom expression handler errors = %v", p.Errors())
		return
	}

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		return
	}

	if program.String() != "handled_expression" {
		t.Errorf("Expected custom expression result, got %v", program.String())
	}
}

func TestParseMultipleStatements(t *testing.T) {
	input := `let x = 5;
	let y = 10;
	x + y`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("Parse errors = %v", p.Errors())
		return
	}

	if len(program.Statements) != 3 {
		t.Errorf("Expected 3 statements, got %d", len(program.Statements))
		return
	}

	expected := "let x=5;let y=10;(x + y)"
	if program.String() != expected {
		t.Errorf("Expected %q, got %q", expected, program.String())
	}
}

func TestParserErrors(t *testing.T) {
	input := "let x = 5;"

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Errorf("Expected no errors for valid input, got %v", p.Errors())
	}

	errors := p.Errors()
	if len(errors) != len(p.errors) {
		t.Errorf("Errors() method returned different count than internal errors")
	}
}

func TestNextToken(t *testing.T) {
	input := "let x = 5"

	l := lexer.New(input)
	p := New(l)

	if p.CurrentToken.Literal != "let" {
		t.Errorf("Expected current token 'let', got %q", p.CurrentToken.Literal)
	}

	if p.PeekToken.Literal != "x" {
		t.Errorf("Expected peek token 'x', got %q", p.PeekToken.Literal)
	}

	p.NextToken()

	if p.CurrentToken.Literal != "x" {
		t.Errorf("After NextToken, expected current token 'x', got %q", p.CurrentToken.Literal)
	}
}
