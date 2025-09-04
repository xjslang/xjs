package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

type TranspilationTest struct {
	name           string
	inputFile      string
	expectedOutput string
}

func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// executeJavaScript executes JavaScript code using Goja and returns the output
func executeJavaScript(code string) (string, error) {
	vm := goja.New()
	var output strings.Builder
	_ = vm.Set("console", map[string]interface{}{
		"log": func(args ...interface{}) {
			for i, arg := range args {
				if i > 0 {
					output.WriteString(" ")
				}
				if arg == nil {
					output.WriteString("null")
				} else {
					output.WriteString(fmt.Sprintf("%v", arg))
				}
			}
			output.WriteString("\n")
		},
	})
	_, err := vm.RunString(code)
	if err != nil {
		return "", fmt.Errorf("failed to execute JavaScript: %v", err)
	}
	result := strings.TrimSpace(output.String())
	return normalizeLineEndings(result), nil
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
		expectedOutput: normalizeLineEndings(strings.TrimSpace(string(expectedOutput))),
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
		actualOutput = normalizeLineEndings(strings.TrimSpace(actualOutput))
		if actualOutput != test.expectedOutput {
			t.Errorf("Output mismatch:\nExpected: %q\nActual:   %q\nTranspiled JS:\n%s",
				test.expectedOutput, actualOutput, transpiledJS)
		}
	})
}

// TestTranspilation tests the transpilation of XJS code to JavaScript by executing it
func TestTranspilation(t *testing.T) {
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

// TestTranspilationErrors tests error handling in transpilation
func TestTranspilationErrors(t *testing.T) {
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
