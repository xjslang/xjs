package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func TestWithComments(t *testing.T) {
	input := `console.log('Hello, World!') // prints a message`
	lb := lexer.NewBuilder()
	p := NewBuilder(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("ParseProgram(%q) error = %v", input, err)
	}
	result := compiler.New().Compile(program)
	fmt.Println(result.Code)
}

func TestParseBasicLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"integer literal", "42", "42;"},
		{"float literal", "3.14", "3.14;"},
		{"string literal", `"hello"`, `"hello";`},
		{"boolean true", "true", "true;"},
		{"boolean false", "false", "false;"},
		{"null literal", "null", "null;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram(%q) error = %v", tt.input, err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"addition", "1 + 2", "(1+2);"},
		{"subtraction", "5 - 3", "(5-3);"},
		{"multiplication", "3 * 4", "(3*4);"},
		{"division", "8 / 2", "(8/2);"},
		{"modulo", "10 % 3", "(10%3);"},
		{"equality", "x == y", "(x==y);"},
		{"inequality", "x != y", "(x!=y);"},
		{"less than", "x < y", "(x<y);"},
		{"greater than", "x > y", "(x>y);"},
		{"logical and", "a && b", "(a&&b);"},
		{"logical or", "a || b", "(a||b);"},
		{"operator precedence", "1 + 2 * 3", "(1+(2*3));"},
		{"parentheses override", "(1 + 2) * 3", "(((1+2))*3);"},
		{"complex precedence", "1 + 2 * 3 - 4", "((1+(2*3))-4);"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram(%q) error = %v", tt.input, err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"negation", "-x", "(-x);"},
		{"logical not", "!true", "(!true);"},
		{"prefix increment", "++x", "(++x);"},
		{"prefix decrement", "--x", "(--x);"},
		{"postfix increment", "x++", "(x++);"},
		{"postfix decrement", "x--", "(x--);"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"simple identifier", "x", "x;"},
		{"camelCase identifier", "myVariable", "myVariable;"},
		{"underscore identifier", "_private", "_private;"},
		{"dollar sign identifier", "$special", "$special;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"let without assignment", "let x", "let x;"},
		{"let with assignment", "let x = 5", "let x=5;"},
		{"let with expression", "let result = 1 + 2", "let result=(1+2);"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"function with body", "function greet() { return 'hello' }", "function greet(){return \"hello\";}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"empty array", "[]", "[];"},
		{"array with numbers", "[1, 2, 3]", "[1,2,3];"},
		{"array with mixed types", "[1, 'hello', true]", "[1,\"hello\",true];"},
		{"nested arrays", "[[1, 2], [3, 4]]", "[[1,2],[3,4]];"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"simple assignment", "x = 5", "x=5;"},
		{"compound addition", "x += 10", "x +=10;"},
		{"compound subtraction", "x -= 5", "x -=5;"},
		{"assignment with expression", "x = 1 + 2", "x=(1+2);"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"function call no args", "foo()", "foo();"},
		{"function call with args", "add(1, 2)", "add(1,2);"},
		{"chained calls", "obj.method()", "obj.method();"},
		{"nested calls", "outer(inner())", "outer(inner());"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
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
		{"dot notation", "obj.prop", "obj.prop;"},
		{"bracket notation", "obj[key]", "obj[key];"},
		{"chained access", "obj.prop.value", "obj.prop.value;"},
		{"mixed access", "obj.prop[0]", "obj.prop[0];"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Errorf("ParseProgram error = %v", err)
				return
			}

			if len(program.Statements) != 1 {
				t.Errorf("Parse(%q) got %d statements, want 1", tt.input, len(program.Statements))
				return
			}

			result := compiler.New().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
			}
		})
	}
}

func TestUseStatementParser(t *testing.T) {
	input := "custom_statement"

	lb := lexer.NewBuilder()
	pb := NewBuilder(lb)

	pb.UseStatementInterceptor(func(p *Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Literal == "custom_statement" {
			return &ast.ExpressionStatement{
				Expression: &ast.Identifier{
					Token: p.CurrentToken,
					Value: "handled_statement",
				},
			}
		}
		return next()
	})

	p := pb.Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("ParseProgram error = %v", err)
		return
	}

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		return
	}

	result := compiler.New().Compile(program)
	if result.Code != "handled_statement;" {
		t.Errorf("Expected custom parser result, got %v", result.Code)
	}
}

func TestUseExpressionHandler(t *testing.T) {
	input := "special_expr"
	lb := lexer.NewBuilder()
	pb := NewBuilder(lb)

	pb.UseExpressionInterceptor(func(p *Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Literal == "special_expr" {
			return p.ParseRemainingExpression(&ast.Identifier{
				Token: p.CurrentToken,
				Value: "handled_expression",
			})
		}
		return next()
	})

	p := pb.Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("ParseProgram error = %v", err)
		return
	}

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		return
	}

	result := compiler.New().Compile(program)
	if result.Code != "handled_expression;" {
		t.Errorf("Expected custom expression result, got %v", result.Code)
	}
}

func TestParseMultipleStatements(t *testing.T) {
	input := `let x = 5;
	let y = 10;
	x + y`

	lb := lexer.NewBuilder()
	p := NewBuilder(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("ParseProgram error = %v", err)
		return
	}

	if len(program.Statements) != 3 {
		t.Errorf("Expected 3 statements, got %d", len(program.Statements))
		return
	}

	expected := "let x=5;let y=10;(x+y);"
	result := compiler.New().Compile(program)
	if result.Code != expected {
		t.Errorf("Expected %q, got %q", expected, result.Code)
	}
}

func TestParserErrors(t *testing.T) {
	input := "let x = 5;"

	lb := lexer.NewBuilder()
	p := NewBuilder(lb).Build(input)
	_, err := p.ParseProgram()
	if err != nil {
		t.Errorf("ParseProgram error = %v", err)
		return
	}
}

func TestNextToken(t *testing.T) {
	input := "let x = 5"

	lb := lexer.NewBuilder()
	p := NewBuilder(lb).Build(input)

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

func TestRegisterOperatorsDuplication(t *testing.T) {
	t.Run("RegisterPrefixOperator duplicates", func(t *testing.T) {
		lb := lexer.NewBuilder()
		pb := NewBuilder(lb)

		// Test builtin prefix operators
		builtinPrefixOps := []token.Type{
			token.IDENT, token.INT, token.FLOAT, token.STRING, token.RAW_STRING,
			token.TRUE, token.FALSE, token.NULL, token.NOT, token.MINUS,
			token.INCREMENT, token.DECREMENT, token.LPAREN, token.LBRACKET,
			token.LBRACE, token.FUNCTION,
		}

		for _, tokenType := range builtinPrefixOps {
			err := pb.RegisterPrefixOperator(tokenType, func(tok token.Token, right func() ast.Expression) ast.Expression {
				return nil
			})
			if err == nil {
				t.Errorf("Expected error when registering duplicate builtin prefix operator %v, but got nil", tokenType)
			} else if !strings.Contains(err.Error(), "duplicate") || !strings.Contains(err.Error(), "prefix") {
				t.Errorf("Expected error message to contain 'duplicate' and 'prefix', but got: %v", err.Error())
			}
		}

		// Test custom prefix operator duplication
		customType := token.Type(1000)
		err := pb.RegisterPrefixOperator(customType, func(tok token.Token, right func() ast.Expression) ast.Expression {
			return nil
		})
		if err != nil {
			t.Errorf("Expected no error when registering new custom prefix operator, but got: %v", err)
		}

		err = pb.RegisterPrefixOperator(customType, func(tok token.Token, right func() ast.Expression) ast.Expression {
			return nil
		})
		if err == nil {
			t.Error("Expected error when registering duplicate custom prefix operator, but got nil")
		}
	})

	t.Run("RegisterInfixOperator duplicates", func(t *testing.T) {
		lb := lexer.NewBuilder()
		pb := NewBuilder(lb)

		// Test builtin infix operators (from precedences map)
		builtinInfixOps := []token.Type{
			token.PLUS, token.MINUS, token.MULTIPLY, token.DIVIDE, token.MODULO,
			token.EQ, token.NOT_EQ, token.LT, token.GT, token.LTE, token.GTE,
			token.AND, token.OR, token.ASSIGN, token.PLUS_ASSIGN, token.MINUS_ASSIGN,
			token.LPAREN, token.DOT, token.LBRACKET, token.INCREMENT, token.DECREMENT,
		}

		for _, tokenType := range builtinInfixOps {
			err := pb.RegisterInfixOperator(tokenType, 1, func(tok token.Token, left ast.Expression, right func() ast.Expression) ast.Expression {
				return nil
			})
			if err == nil {
				t.Errorf("Expected error when registering duplicate builtin infix operator %v, but got nil", tokenType)
			} else if !strings.Contains(err.Error(), "duplicate") || !strings.Contains(err.Error(), "infix") {
				t.Errorf("Expected error message to contain 'duplicate' and 'infix', but got: %v", err.Error())
			}
		}

		// Test custom infix operator duplication
		customType := token.Type(1001)
		err := pb.RegisterInfixOperator(customType, 5, func(tok token.Token, left ast.Expression, right func() ast.Expression) ast.Expression {
			return nil
		})
		if err != nil {
			t.Errorf("Expected no error when registering new custom infix operator, but got: %v", err)
		}

		err = pb.RegisterInfixOperator(customType, 5, func(tok token.Token, left ast.Expression, right func() ast.Expression) ast.Expression {
			return nil
		})
		if err == nil {
			t.Error("Expected error when registering duplicate custom infix operator, but got nil")
		}
	})

	t.Run("RegisterPostfixOperator duplicates", func(t *testing.T) {
		lb := lexer.NewBuilder()
		pb := NewBuilder(lb)

		// Test builtin postfix operators
		builtinPostfixOps := []token.Type{
			token.INCREMENT, token.DECREMENT,
		}

		for _, tokenType := range builtinPostfixOps {
			err := pb.RegisterPostfixOperator(tokenType, func(tok token.Token, left ast.Expression) ast.Expression {
				return nil
			})
			if err == nil {
				t.Errorf("Expected error when registering duplicate builtin postfix operator %v, but got nil", tokenType)
			} else if !strings.Contains(err.Error(), "duplicate") || !strings.Contains(err.Error(), "postfix") {
				t.Errorf("Expected error message to contain 'duplicate' and 'postfix', but got: %v", err.Error())
			}
		}

		// Test custom postfix operator duplication
		customType := token.Type(1002)
		err := pb.RegisterPostfixOperator(customType, func(tok token.Token, left ast.Expression) ast.Expression {
			return nil
		})
		if err != nil {
			t.Errorf("Expected no error when registering new custom postfix operator, but got: %v", err)
		}

		err = pb.RegisterPostfixOperator(customType, func(tok token.Token, left ast.Expression) ast.Expression {
			return nil
		})
		if err == nil {
			t.Error("Expected error when registering duplicate custom postfix operator, but got nil")
		}
	})
}

func TestWithSmartSemicolon(t *testing.T) {
	tests := []struct {
		name     string
		smart    bool
		input    string
		expected string
	}{
		{"with smart semicolons", true, "console.log('first line')\n(function() { console.log('second line') })()", "console.log(\"first line\");(function(){console.log(\"second line\");})();"},
		{"without smart semicolons", false, "console.log('first line')\n(function() { console.log('second line') })()", "console.log(\"first line\")(function(){console.log(\"second line\");})();"},
	}
	for _, tt := range tests {
		lb := lexer.NewBuilder()
		p := NewBuilder(lb).
			WithSmartSemicolon(tt.smart).
			Build(tt.input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Errorf("ParseProgram error = %v", err)
			return
		}

		result := compiler.New().Compile(program)
		if result.Code != tt.expected {
			t.Errorf("Parse(%q) = %v, want %v", tt.input, result.Code, tt.expected)
		}
	}
}
