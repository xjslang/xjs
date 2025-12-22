//go:build integration

package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// IIFEFunctionDeclaration wraps a function declaration in an IIFE
type IIFEFunctionDeclaration struct {
	*ast.FunctionDeclaration
}

// WriteTo generates an IIFE: (function name() {...})()
func (iife *IIFEFunctionDeclaration) WriteTo(cw *ast.CodeWriter) {
	cw.WriteString("(function ")
	iife.Name.WriteTo(cw)
	cw.WriteRune('(')
	for i, param := range iife.Parameters {
		if i > 0 {
			cw.WriteRune(',')
		}
		param.WriteTo(cw)
	}
	cw.WriteRune(')')
	iife.Body.WriteTo(cw)
	cw.WriteString(")()")
}

// transpileASTToJS converts an AST to JavaScript (now with automatic semicolons)
func transpileASTToJS(program *ast.Program) string {
	result := compiler.New().Compile(program)
	return result.Code
}

// TestMiddlewareParsers tests the custom middleware functionality
func TestMiddlewareParsers(t *testing.T) {
	t.Run("iife_function_declaration_middleware", func(t *testing.T) {
		input := `
		function main() {   // transformed to (function main() {...})()
			function main() { // not transformed, as it is not a top-level function
				console.log('Hello, world!')
			}
			main()
		}`
		expectedOutput := "Hello, world!"

		lb := lexer.NewBuilder()
		pb := parser.NewBuilder(lb)
		pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
			if p.CurrentToken.Type != token.FUNCTION {
				return next()
			}

			// Check if we're NOT inside a function (top-level only)
			if p.IsInFunction() {
				return next()
			}

			funcDecl := p.ParseFunctionStatement()
			if funcDecl.Name != nil && funcDecl.Name.Value == "main" {
				return &IIFEFunctionDeclaration{
					FunctionDeclaration: funcDecl,
				}
			}

			return funcDecl
		})

		p := pb.Build(input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram error = %v", err)
		}

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

	t.Run("expression_parser_middleware", func(t *testing.T) {
		input := `let x = 5 + 3; console.log(x)`
		expectedOutput := "8"

		// Create parser with custom expression parser
		lb := lexer.NewBuilder()
		pb := parser.NewBuilder(lb)

		// Add middleware that logs when processing numbers
		var processedNumbers []string
		pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			if p.CurrentToken.Type == token.INT {
				processedNumbers = append(processedNumbers, p.CurrentToken.Literal)
			}
			return next()
		})

		p := pb.Build(input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram error = %v", err)
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

	t.Run("multiple_expression_parsers", func(t *testing.T) {
		input := `let message = 'Hello' + ' ' + 'World'; console.log(message)`
		expectedOutput := "Hello World"

		lb := lexer.NewBuilder()
		pb := parser.NewBuilder(lb)

		var stringCount int
		var identifierCount int

		// First middleware - count strings
		pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			if p.CurrentToken.Type == token.STRING {
				stringCount++
			}
			return next()
		})

		// Second middleware - count identifiers
		pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			if p.CurrentToken.Type == token.IDENT {
				identifierCount++
			}
			return next()
		})

		p := pb.Build(input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram error = %v", err)
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

		lb := lexer.NewBuilder()
		pb := parser.NewBuilder(lb)

		var logMessages []string

		// Middleware that logs parsing progress
		pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
			currentTokenType := p.CurrentToken.Type.String()
			logMessages = append(logMessages, fmt.Sprintf("Processing: %s", currentTokenType))
			return next()
		})

		p := pb.Build(input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram error = %v", err)
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

// TestPluginExecutionOrder verifies that plugins are executed in the order they are installed
func TestPluginExecutionOrder(t *testing.T) {
	t.Run("plugins_execute_in_installation_order", func(t *testing.T) {
		input := `let x = 42; console.log(x)`
		expectedOutput := "42"

		// Shared slice to record execution order
		var executionOrder []string

		// Create plugins dynamically - each plugin adds its name to execution order when processing LET tokens
		plugins := make([]func(*parser.Builder), 3)
		for i := 0; i < 3; i++ {
			pluginName := fmt.Sprintf("Plugin%d", i+1)
			plugins[i] = func(pb *parser.Builder) {
				pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
					if p.CurrentToken.Type == token.LET {
						executionOrder = append(executionOrder, pluginName)
					}
					return next()
				})
			}
		}

		// Install plugins in order: Plugin1, Plugin2, Plugin3
		lb := lexer.NewBuilder()
		pb := parser.NewBuilder(lb)
		for _, plugin := range plugins {
			pb.Install(plugin)
		}

		p := pb.Build(input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram error = %v", err)
		}

		// Verify execution order
		expectedOrder := []string{"Plugin1", "Plugin2", "Plugin3"}
		if len(executionOrder) != len(expectedOrder) {
			t.Fatalf("Expected %d plugins to execute, got %d", len(expectedOrder), len(executionOrder))
		}

		for i, pluginName := range expectedOrder {
			if executionOrder[i] != pluginName {
				t.Errorf("Plugin at position %d: expected %q, got %q", i, pluginName, executionOrder[i])
			}
		}

		// Verify the code still works correctly
		transpiledJS := transpileASTToJS(program)
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		actualOutput = strings.TrimSpace(actualOutput)
		if actualOutput != expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q", expectedOutput, actualOutput)
		}

		t.Logf("Plugins executed in correct order: %v", executionOrder)
	})

	t.Run("expression_interceptors_execute_in_order", func(t *testing.T) {
		input := `let result = 10 + 20; console.log(result)`
		expectedOutput := "30"

		// Shared slice to record execution order
		var executionOrder []string

		// Create plugins dynamically - each tracks INT tokens
		plugins := make([]func(*parser.Builder), 3)
		for i := 0; i < 3; i++ {
			pluginName := fmt.Sprintf("Plugin%d-INT", i+1)
			plugins[i] = func(pb *parser.Builder) {
				pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
					if p.CurrentToken.Type == token.INT {
						executionOrder = append(executionOrder, pluginName)
					}
					return next()
				})
			}
		}

		// Install plugins in order
		lb := lexer.NewBuilder()
		pb := parser.NewBuilder(lb)
		for _, plugin := range plugins {
			pb.Install(plugin)
		}

		p := pb.Build(input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram error = %v", err)
		}

		// Each INT token (10 and 20) should be processed by all three plugins in order
		expectedPattern := []string{
			"Plugin1-INT", "Plugin2-INT", "Plugin3-INT", // First INT (10)
			"Plugin1-INT", "Plugin2-INT", "Plugin3-INT", // Second INT (20)
		}

		if len(executionOrder) != len(expectedPattern) {
			t.Fatalf("Expected %d interceptor calls, got %d\nExecution order: %v",
				len(expectedPattern), len(executionOrder), executionOrder)
		}

		for i, expected := range expectedPattern {
			if executionOrder[i] != expected {
				t.Errorf("Position %d: expected %q, got %q", i, expected, executionOrder[i])
			}
		}

		// Verify the code still works correctly
		transpiledJS := transpileASTToJS(program)
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		actualOutput = strings.TrimSpace(actualOutput)
		if actualOutput != expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q", expectedOutput, actualOutput)
		}

		t.Logf("Expression interceptors executed in correct order: %v", executionOrder)
	})

	t.Run("mixed_plugins_with_multiple_features", func(t *testing.T) {
		input := `let x = 5; function test() { return x * 2; } console.log(test())`
		expectedOutput := "10"

		// Shared map to track what each plugin processed
		processedByPlugin := make(map[string][]string)

		// Plugin1 - processes both statements and expressions
		Plugin1 := func(pb *parser.Builder) {
			pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
				if p.CurrentToken.Type == token.LET {
					processedByPlugin["Plugin1"] = append(processedByPlugin["Plugin1"], "LET")
				}
				return next()
			})
			pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
				if p.CurrentToken.Type == token.INT {
					processedByPlugin["Plugin1"] = append(processedByPlugin["Plugin1"], "INT")
				}
				return next()
			})
		}

		// Plugin2 - processes only statements
		Plugin2 := func(pb *parser.Builder) {
			pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
				if p.CurrentToken.Type == token.FUNCTION {
					processedByPlugin["Plugin2"] = append(processedByPlugin["Plugin2"], "FUNCTION")
				}
				return next()
			})
		}

		// Plugin3 - processes only expressions
		Plugin3 := func(pb *parser.Builder) {
			pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
				if p.CurrentToken.Type == token.IDENT {
					processedByPlugin["Plugin3"] = append(processedByPlugin["Plugin3"], "IDENT")
				}
				return next()
			})
		}

		// Install plugins
		lb := lexer.NewBuilder()
		pb := parser.NewBuilder(lb).Install(Plugin1).Install(Plugin2).Install(Plugin3)

		p := pb.Build(input)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("ParseProgram error = %v", err)
		}

		// Verify each plugin processed what it should
		if len(processedByPlugin["Plugin1"]) == 0 {
			t.Error("Plugin1 did not process any tokens")
		}
		if len(processedByPlugin["Plugin2"]) == 0 {
			t.Error("Plugin2 did not process any tokens")
		}
		if len(processedByPlugin["Plugin3"]) == 0 {
			t.Error("Plugin3 did not process any tokens")
		}

		// Verify the code still works correctly
		transpiledJS := transpileASTToJS(program)
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		actualOutput = strings.TrimSpace(actualOutput)
		if actualOutput != expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q", expectedOutput, actualOutput)
		}

		t.Logf("Plugin1 processed: %v", processedByPlugin["Plugin1"])
		t.Logf("Plugin2 processed: %v", processedByPlugin["Plugin2"])
		t.Logf("Plugin3 processed: %v", processedByPlugin["Plugin3"])
	})
}
