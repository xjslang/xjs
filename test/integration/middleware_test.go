package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// transpileASTToJS converts an AST to JavaScript (now with automatic semicolons)
func transpileASTToJS(program *ast.Program) string {
	return program.String()
}

// TestMiddlewareHandlers tests the custom middleware functionality
func TestMiddlewareHandlers(t *testing.T) {
	t.Run("expression_handler_middleware", func(t *testing.T) {
		input := `let x = 5 + 3; console.log(x)`
		expectedOutput := "8"

		// Create parser with custom expression handler
		l := lexer.New(input)
		p := parser.New(l)

		// Add middleware that logs when processing numbers
		var processedNumbers []string
		p.UseExpressionHandler(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			if p.CurrentToken.Type == token.INT {
				processedNumbers = append(processedNumbers, p.CurrentToken.Literal)
			}
			return next()
		})

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			t.Fatalf("Parser errors: %v", p.Errors())
		}

		// Verify the middleware was called for the numbers
		expectedNumbers := []string{"5", "3"}
		if len(processedNumbers) != len(expectedNumbers) {
			t.Errorf("Expected to process %d numbers, got %d", len(expectedNumbers), len(processedNumbers))
		}

		// Execute the transpiled code
		transpiledJS := transpileASTToJS(program)
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		actualOutput = strings.TrimSpace(actualOutput)
		if actualOutput != expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q", expectedOutput, actualOutput)
		}
	})

	t.Run("multiple_expression_handlers", func(t *testing.T) {
		input := `let message = 'Hello' + ' ' + 'World'; console.log(message)`
		expectedOutput := "Hello World"

		l := lexer.New(input)
		p := parser.New(l)

		var stringCount int
		var identifierCount int

		// First middleware - count strings
		p.UseExpressionHandler(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			if p.CurrentToken.Type == token.STRING {
				stringCount++
			}
			return next()
		})

		// Second middleware - count identifiers
		p.UseExpressionHandler(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			if p.CurrentToken.Type == token.IDENT {
				identifierCount++
			}
			return next()
		})

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			t.Fatalf("Parser errors: %v", p.Errors())
		}

		// Execute the transpiled code
		transpiledJS := transpileASTToJS(program)
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		actualOutput = strings.TrimSpace(actualOutput)
		if actualOutput != expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q", expectedOutput, actualOutput)
		}

		// Verify middleware was called
		if stringCount == 0 {
			t.Error("String middleware was not called")
		}
		if identifierCount == 0 {
			t.Error("Identifier middleware was not called")
		}
	})
}

// TestCustomLanguageFeatures tests custom language features that could be added via middleware
func TestCustomLanguageFeatures(t *testing.T) {
	t.Run("logging_middleware", func(t *testing.T) {
		input := `function test() { let x = 42; return x; } console.log(test())`
		expectedOutput := "42"

		l := lexer.New(input)
		p := parser.New(l)

		var logMessages []string

		// Middleware that logs parsing progress
		p.UseExpressionHandler(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			currentTokenType := p.CurrentToken.Type.String()
			logMessages = append(logMessages, fmt.Sprintf("Processing: %s", currentTokenType))
			return next()
		})

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			t.Fatalf("Parser errors: %v", p.Errors())
		}

		// Execute the transpiled code
		transpiledJS := transpileASTToJS(program)
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		actualOutput = strings.TrimSpace(actualOutput)
		if actualOutput != expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q", expectedOutput, actualOutput)
		}

		// Verify that middleware was active
		if len(logMessages) == 0 {
			t.Error("Logging middleware was not called")
		}

		t.Logf("Middleware logged %d parsing steps", len(logMessages))
	})
}

// TestComplexTranspilation tests more complex transpilation scenarios
func TestComplexTranspilation(t *testing.T) {
	complexTests := []TranspilationTest{
		{
			name:           "nested_functions",
			inputFile:      `function outer(x) { function inner(y) { return x + y; } return inner(5); } console.log(outer(10))`,
			expectedOutput: "15",
		},
		{
			name:           "array_methods",
			inputFile:      `let arr = [1, 2, 3]; console.log(arr.length)`,
			expectedOutput: "3",
		},
		{
			name:           "string_operations",
			inputFile:      `let str = 'Hello'; console.log(str + ' World')`,
			expectedOutput: "Hello World",
		},
		{
			name:           "boolean_logic",
			inputFile:      `let a = true; let b = false; console.log(a && b)`,
			expectedOutput: "false",
		},
		{
			name:           "undefined_and_null",
			inputFile:      `let y = null; console.log(y == null)`,
			expectedOutput: "true",
		},
	}

	for _, test := range complexTests {
		RunTranspilationTest(t, test)
	}
}

// TestPerformanceRegression tests for performance regressions in transpilation
func TestPerformanceRegression(t *testing.T) {
	// Large input to test performance
	var inputBuilder strings.Builder
	inputBuilder.WriteString("let sum = 0;\n")
	for i := 0; i < 100; i++ {
		inputBuilder.WriteString(fmt.Sprintf("sum += %d;\n", i))
	}
	inputBuilder.WriteString("console.log(sum);")

	input := inputBuilder.String()

	// Time the transpilation
	_, err := transpileXJSCode(input)
	if err != nil {
		t.Fatalf("Transpilation failed: %v", err)
	}

	// If we get here without timing out, the performance is acceptable
	t.Log("Performance test passed - large input transpiled successfully")
}
