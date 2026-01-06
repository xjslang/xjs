package integration

import (
	"strings"
	"testing"

	"github.com/xjslang/xjs/compiler"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestPrettyPrint_LineComments(t *testing.T) {
	input := `
// Configuration file
let host = 'localhost'
let port = 3000

// Main function
function main() {
	// Initialize server
	console.log('Starting server')
}`

	expected := `// Configuration file
let host = "localhost";
let port = 3000;

// Main function
function main() {
  // Initialize server
  console.log("Starting server");
}`

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := compiler.New().WithPrettyPrint().Compile(program)
	if result.Code != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result.Code)
	}
}

func TestPrettyPrint_InlineComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "inline comments on let statements",
			input: `
let x = 5 // without semicolon
let y = 10; // with semicolon`,
			expected: `let x = 5; // without semicolon
let y = 10; // with semicolon`,
		},
		{
			name: "inline comments on return statements",
			input: `
function test() {
	return x // without semicolon
}
function test2() {
	return y; // with semicolon
}`,
			expected: `function test() {
  return x; // without semicolon
}
function test2() {
  return y; // with semicolon
}`,
		},
		{
			name: "inline comments on expression statements",
			input: `
console.log('Hello') // message
console.log('World'); // another message`,
			expected: `console.log("Hello"); // message
console.log("World"); // another message`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := compiler.New().WithPrettyPrint().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result.Code)
			}
		})
	}
}

func TestPrettyPrint_ObjectLiteralComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "object with inline comments",
			input: `
let user = {
	name: 'John Smith', // user name
	age: 35 // user age
}`,
			expected: `let user = {
  name: "John Smith", // user name
  age: 35 // user age
};`,
		},
		{
			name: "object with inline comments and trailing comma",
			input: `
let config = {
	host: 'localhost', // server host
	port: 3000, // server port
}`,
			expected: `let config = {
  host: "localhost", // server host
  port: 3000 // server port
};`,
		},
		{
			name: "object without comments (single line)",
			input: `
let simple = {name: 'John', age: 30}`,
			expected: `let simple = {name: "John", age: 30};`,
		},
		{
			name: "nested objects with comments",
			input: `
let server = {
	config: {
		host: 'localhost', // server host
		port: 8080 // server port
	}, // configuration object
	name: 'MyServer' // server name
}`,
			expected: `let server = {
  config: {
    host: "localhost", // server host
    port: 8080 // server port
  }, // configuration object
  name: "MyServer" // server name
};`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).Build(tt.input)
			program, err := p.ParseProgram()
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			result := compiler.New().WithPrettyPrint().Compile(program)
			if result.Code != tt.expected {
				t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result.Code)
			}
		})
	}
}

func TestPrettyPrint_MixedComments(t *testing.T) {
	input := `
// Application configuration
let config = {
	host: 'localhost', // default host
	port: 3000, // default port
}

// Initialize application
function init() {
	let message = 'Starting' // status message
	// Log to console
	console.log(message, config.host) // debug output
	return config // return configuration
}`

	expected := `// Application configuration
let config = {
  host: "localhost", // default host
  port: 3000 // default port
};

// Initialize application
function init() {
  let message = "Starting"; // status message
  // Log to console
  console.log(message, config.host); // debug output
  return config; // return configuration
}`

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	result := compiler.New().WithPrettyPrint().Compile(program)

	// Normalize line endings for comparison
	expectedNormalized := strings.ReplaceAll(expected, "\r\n", "\n")
	resultNormalized := strings.ReplaceAll(result.Code, "\r\n", "\n")

	if resultNormalized != expectedNormalized {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, result.Code)
	}
}
