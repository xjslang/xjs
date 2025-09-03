package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

type TranspilationTest struct {
	name           string
	inputFile      string
	expectedOutput string
}

// executeJavaScript executes JavaScript code using Node.js and returns the output
func executeJavaScript(code string) (string, error) {
	// Create a temporary file with the transpiled JavaScript
	tempFile, err := os.CreateTemp("", "xjs_test_*.js")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()

	// Write the transpiled code to the temporary file
	_, err = tempFile.WriteString(code)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}
	_ = tempFile.Close()

	// Execute the JavaScript file using Node.js
	cmd := exec.Command("node", tempFile.Name())
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("execution failed with exit code %d: %s", exitError.ExitCode(), string(exitError.Stderr))
		}
		return "", fmt.Errorf("failed to execute JavaScript: %v", err)
	}

	return string(output), nil
}

// checkNodeJSAvailability verifies that Node.js is available for testing
func checkNodeJSAvailability(t *testing.T) {
	_, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node.js is not available, skipping transpilation tests")
	}
}

// transpileXJSCode transpiles XJS code to JavaScript using the main Parse function
func transpileXJSCode(input string) (string, error) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) > 0 {
		return "", fmt.Errorf("parser errors: %v", errors)
	}

	// Convert the AST to JavaScript code (now with automatic semicolons)
	result := program.String()

	return result, nil
}

// loadTestCase loads a test case from fixture files
func loadTestCase(t *testing.T, baseName string) TranspilationTest {
	inputFile := filepath.Join("../testdata", baseName+".js")
	outputFile := filepath.Join("../testdata", baseName+".output")

	// Read input file
	inputContent, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Failed to read input file %s: %v", inputFile, err)
	}

	// Read expected output file
	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file %s: %v", outputFile, err)
	}

	return TranspilationTest{
		name:           baseName,
		inputFile:      string(inputContent),
		expectedOutput: strings.TrimSpace(string(expectedOutput)),
	}
}

// RunTranspilationTest runs a single transpilation test
func RunTranspilationTest(t *testing.T, test TranspilationTest) {
	t.Run(test.name, func(t *testing.T) {
		// Transpile the XJS code to JavaScript
		transpiledJS, err := transpileXJSCode(test.inputFile)
		if err != nil {
			t.Fatalf("Transpilation failed: %v", err)
		}

		// Execute the transpiled JavaScript
		actualOutput, err := executeJavaScript(transpiledJS)
		if err != nil {
			t.Fatalf("JavaScript execution failed: %v", err)
		}

		// Compare the actual output with expected output
		actualOutput = strings.TrimSpace(actualOutput)
		if actualOutput != test.expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q\nTranspiled JS:\n%s",
				test.expectedOutput, actualOutput, transpiledJS)
		}
	})
}

// TestTranspilation tests the transpilation of XJS code to JavaScript by executing it
func TestTranspilation(t *testing.T) {
	checkNodeJSAvailability(t)

	// Test cases based on fixture files
	testCases := []string{
		"basic",
		"function",
		"array_loop",
		"conditional",
		"object",
		"operators",
		"loops",
		"data_types",
		"assignment",
		"compound_assignment",
		"function_expressions",
		"property_access",
		"complex_conditionals",
		"complex_expressions",
		"type_coercion",
		"strict_equality",
		"comments",
		"comments_comprehensive",
		"increment_decrement",
	}

	for _, testCase := range testCases {
		test := loadTestCase(t, testCase)
		RunTranspilationTest(t, test)
	}
}

// TestTranspilationBasicInline tests basic transpilation with inline test data
func TestTranspilationBasicInline(t *testing.T) {
	checkNodeJSAvailability(t)

	tests := []TranspilationTest{
		{
			name:           "simple_console_log",
			inputFile:      `console.log('Hello, World!')`,
			expectedOutput: "Hello, World!",
		},
		{
			name:           "variable_declaration_and_log",
			inputFile:      `let message = 'Test Message'; console.log(message)`,
			expectedOutput: "Test Message",
		},
		{
			name:           "simple_arithmetic",
			inputFile:      `let result = 5 + 3; console.log(result)`,
			expectedOutput: "8",
		},
		{
			name:           "function_declaration_and_call",
			inputFile:      `function greet(name) { return 'Hello, ' + name; } console.log(greet('XJS'))`,
			expectedOutput: "Hello, XJS",
		},
	}

	for _, test := range tests {
		RunTranspilationTest(t, test)
	}
}

// TestTranspilationErrors tests error handling in transpilation
func TestTranspilationErrors(t *testing.T) {
	checkNodeJSAvailability(t)

	errorTests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid_syntax",
			input: "let x = ",
		},
		{
			name:  "unclosed_string",
			input: `console.log("unclosed string`,
		},
	}

	for _, test := range errorTests {
		t.Run(test.name, func(t *testing.T) {
			_, err := transpileXJSCode(test.input)
			if err == nil {
				t.Errorf("Expected transpilation error for input: %s", test.input)
			}
		})
	}
}

// BenchmarkTranspilation benchmarks the transpilation process
func BenchmarkTranspilation(b *testing.B) {
	input := `
		function fibonacci(n) {
			if (n <= 1) return n;
			return fibonacci(n - 1) + fibonacci(n - 2);
		}
		console.log(fibonacci(10));
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := transpileXJSCode(input)
		if err != nil {
			b.Fatalf("Transpilation failed: %v", err)
		}
	}
}
