//go:build integration

package integration

import (
	"testing"

	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestTolerantMode(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		strictShouldPass    bool
		tolerantShouldParse bool
	}{
		{
			name:                "valid code with semicolons",
			input:               "let x = 42;",
			strictShouldPass:    true,
			tolerantShouldParse: true,
		},
		{
			name:                "valid code without semicolons (ASI)",
			input:               "let x = 42\nlet y = 10",
			strictShouldPass:    true,
			tolerantShouldParse: true,
		},
		{
			name:                "invalid code - missing semicolon on same line",
			input:               "let x = 1 let y = 2",
			strictShouldPass:    false,
			tolerantShouldParse: true,
		},
		{
			name:                "incomplete expression",
			input:               "let x = ",
			strictShouldPass:    false,
			tolerantShouldParse: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test strict mode
			t.Run("strict", func(t *testing.T) {
				lb := lexer.NewBuilder()
				p := parser.NewBuilder(lb).Build(tt.input)
				_, err := p.ParseProgram()

				if tt.strictShouldPass && err != nil {
					t.Errorf("strict mode: expected no error, got %v", err)
				}
				if !tt.strictShouldPass && err == nil {
					t.Errorf("strict mode: expected error, got none")
				}
			})

			// Test tolerant mode
			t.Run("tolerant", func(t *testing.T) {
				lb := lexer.NewBuilder()
				p := parser.NewBuilder(lb).WithTolerantMode(true).Build(tt.input)
				program, _ := p.ParseProgram()

				if tt.tolerantShouldParse && program == nil {
					t.Errorf("tolerant mode: expected to parse something, got nil")
				}
			})
		})
	}
}

func TestTolerantModeWithBuilder(t *testing.T) {
	input := "let x = 1 let y = 2" // Invalid: missing semicolon

	// Strict mode via builder (default)
	t.Run("builder_strict", func(t *testing.T) {
		p := parser.NewBuilder(lexer.NewBuilder()).Build(input)
		_, err := p.ParseProgram()
		if err == nil {
			t.Error("expected error in strict mode, got none")
		}
	})

	// Tolerant mode via builder
	t.Run("builder_tolerant", func(t *testing.T) {
		p := parser.NewBuilder(lexer.NewBuilder()).
			WithTolerantMode(true).
			Build(input)
		program, _ := p.ParseProgram()
		if program == nil {
			t.Fatal("expected program in tolerant mode, got nil")
		}
		if len(program.Statements) == 0 {
			t.Error("expected statements in tolerant mode, got none")
		}
	})
}

func TestTolerantModeContinuesParsing(t *testing.T) {
	// Code with multiple errors
	input := `
		let a = 1 let b = 2
		let c = 3
		let d = 4 let e = 5
	`

	// Strict mode should stop early
	t.Run("strict_stops_early", func(t *testing.T) {
		lb := lexer.NewBuilder()
		p := parser.NewBuilder(lb).Build(input)
		_, err := p.ParseProgram()
		if err == nil {
			t.Error("expected error in strict mode")
		}
	})

	// Tolerant mode should continue and parse valid statements
	t.Run("tolerant_continues", func(t *testing.T) {
		lb := lexer.NewBuilder()
		p := parser.NewBuilder(lb).WithTolerantMode(true).Build(input)
		program, _ := p.ParseProgram()

		if program == nil {
			t.Fatal("expected program in tolerant mode")
		}

		// Should parse at least the valid statement (let c = 3)
		if len(program.Statements) < 1 {
			t.Errorf("expected at least 1 statement, got %d", len(program.Statements))
		}
	})
}

func TestTolerantModeForLSPScenarios(t *testing.T) {
	scenarios := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "incomplete_function",
			input:       "function foo() { let x = ",
			description: "User is typing a function",
		},
		{
			name:        "missing_closing_brace",
			input:       "function foo() { return 42",
			description: "User hasn't closed the function yet",
		},
		{
			name:        "multiple_statements_no_semicolons",
			input:       "let a = 1\nlet b = 2\nlet c",
			description: "User is typing multiple statements",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			lb := lexer.NewBuilder()
			p := parser.NewBuilder(lb).WithTolerantMode(true).Build(scenario.input)
			program, _ := p.ParseProgram()

			// In tolerant mode, we should get some program structure
			// even if the code is incomplete or invalid
			if program == nil {
				t.Fatalf("%s: expected program despite errors", scenario.description)
			}

			t.Logf("%s: parsed %d statements", scenario.description, len(program.Statements))
		})
	}
}
