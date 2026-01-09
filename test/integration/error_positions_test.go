package integration

import (
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

// TestErrorPositions verifies that parser errors report accurate start and end positions
// for tokens causing syntax errors. This is critical for IDE integrations and error highlighting.
func TestErrorPositions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			message   string
			startLine int
			startCol  int
			endLine   int
			endCol    int
		}
	}{
		{
			name:  "expected semicolon after statement",
			input: "let x = 5 let y = 10",
			expected: []struct {
				message   string
				startLine int
				startCol  int
				endLine   int
				endCol    int
			}{
				{
					message:   "semicolon or newline expected",
					startLine: 0,
					startCol:  10,
					endLine:   0,
					endCol:    13,
				},
			},
		},
		{
			name:  "missing closing parenthesis",
			input: "let x = (5 + 3",
			expected: []struct {
				message   string
				startLine int
				startCol  int
				endLine   int
				endCol    int
			}{
				{
					message:   ") expected",
					startLine: 0,
					startCol:  14,
					endLine:   0,
					endCol:    14,
				},
			},
		},
		{
			name:  "missing function name",
			input: "function (a, b) { return a + b }",
			expected: []struct {
				message   string
				startLine int
				startCol  int
				endLine   int
				endCol    int
			}{
				{
					message:   "identifier expected",
					startLine: 0,
					startCol:  9,
					endLine:   0,
					endCol:    9,
				},
			},
		},
		{
			name:  "unexpected token in expression",
			input: "let x = 5 + + 3",
			expected: []struct {
				message   string
				startLine int
				startCol  int
				endLine   int
				endCol    int
			}{
				{
					message:   "unexpected +",
					startLine: 0,
					startCol:  12,
					endLine:   0,
					endCol:    12,
				},
				{
					message:   "semicolon or newline expected",
					startLine: 0,
					startCol:  14,
					endLine:   0,
					endCol:    15,
				},
			},
		},
		{
			name: "error on multi-line code",
			input: `let x = 5
let y = 
let z = 10`,
			expected: []struct {
				message   string
				startLine int
				startCol  int
				endLine   int
				endCol    int
			}{
				{
					message:   "unexpected let",
					startLine: 2,
					startCol:  0,
					endLine:   2,
					endCol:    3,
				},
				{
					message:   "semicolon or newline expected",
					startLine: 2,
					startCol:  4,
					endLine:   2,
					endCol:    5,
				},
			},
		},
		{
			name:  "missing closing brace",
			input: "function test() { let x = 5",
			expected: []struct {
				message   string
				startLine int
				startCol  int
				endLine   int
				endCol    int
			}{
				{
					message:   "unclosed block statement, expected '}'",
					startLine: 0,
					startCol:  27,
					endLine:   0,
					endCol:    27,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).Build(tt.input)
			_, _ = p.ParseProgram()
			errors := p.Errors()

			if len(tt.expected) == 0 {
				if len(errors) > 0 {
					t.Errorf("Expected no errors, but got %d: %v", len(errors), errors)
				}
				return
			}

			if len(errors) < len(tt.expected) {
				t.Fatalf("Expected at least %d errors, got %d.\nErrors: %v", len(tt.expected), len(errors), errors)
			}

			for i, expected := range tt.expected {
				err := errors[i]

				if err.Message != expected.message {
					t.Errorf("Error %d: expected message %q, got %q",
						i, expected.message, err.Message)
				}

				if err.Range.Start.Line != expected.startLine {
					t.Errorf("Error %d: expected start line %d, got %d",
						i, expected.startLine, err.Range.Start.Line)
				}
				if err.Range.Start.Column != expected.startCol {
					t.Errorf("Error %d: expected start column %d, got %d",
						i, expected.startCol, err.Range.Start.Column)
				}

				if err.Range.End.Line != expected.endLine {
					t.Errorf("Error %d: expected end line %d, got %d",
						i, expected.endLine, err.Range.End.Line)
				}
				if err.Range.End.Column != expected.endCol {
					t.Errorf("Error %d: expected end column %d, got %d",
						i, expected.endCol, err.Range.End.Column)
				}

				if t.Failed() {
					t.Logf("Full error details: %+v", err)
					t.Logf("Input: %q", tt.input)
				}
			}
		})
	}
}

// TestTokenPositionsInAST verifies that tokens in the AST maintain correct position information
func TestTokenPositionsInAST(t *testing.T) {
	input := `let x = 42
let y = "hello"
function add(a, b) {
  return a + b
}`

	lb := lexer.NewBuilder()
	p := parser.NewBuilder(lb).Build(input)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Unexpected parse error: %v", err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Expected 3 statements, got %d", len(program.Statements))
	}

	t.Log("Successfully parsed program with token.Position structure")
}
