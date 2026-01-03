//go:build integration

package integration

import (
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/sourcemap"
	"github.com/xjslang/xjs/token"
)

func TestWriteToWithSourceMapper(t *testing.T) {
	// Create a simple AST
	stmt := &ast.LetStatement{
		Token: token.Token{Type: token.LET, Literal: "let", Start: token.Position{Line: 1, Column: 1}, End: token.Position{Line: 1, Column: 1}},
		Name: &ast.Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "x", Start: token.Position{Line: 1, Column: 5}, End: token.Position{Line: 1, Column: 5}},
			Value: "x",
		},
		Value: &ast.IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: "42", Start: token.Position{Line: 1, Column: 9}, End: token.Position{Line: 1, Column: 9}},
		},
	}

	// Generate without source map
	var w1 ast.CodeWriter
	stmt.WriteTo(&w1)
	if w1.String() != "let x=42;" {
		t.Errorf("WriteTo without mapper = %q, want %q", w1.String(), "let x=42")
	}

	// Generate with source map
	w2 := ast.CodeWriter{
		Builder: strings.Builder{},
		Mapper:  sourcemap.New(),
	}
	stmt.WriteTo(&w2)

	if w2.String() != "let x=42;" {
		t.Errorf("WriteTo with mapper = %q, want %q", w2.String(), "let x=42")
	}

	sm := w2.Mapper.SourceMap()
	if sm.Mappings == "" {
		t.Error("Source map mappings should not be empty")
	}

	// Should have captured the identifier name
	if len(sm.Names) != 1 || sm.Names[0] != "x" {
		t.Errorf("Names = %v, want [x]", sm.Names)
	}
}

func TestProgramWithSourceMapper(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let", Start: token.Position{Line: 1, Column: 1}, End: token.Position{Line: 1, Column: 1}},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x", Start: token.Position{Line: 1, Column: 5}, End: token.Position{Line: 1, Column: 5}},
					Value: "x",
				},
				Value: &ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5", Start: token.Position{Line: 1, Column: 9}, End: token.Position{Line: 1, Column: 9}},
				},
			},
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let", Start: token.Position{Line: 2, Column: 1}, End: token.Position{Line: 2, Column: 1}},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "y", Start: token.Position{Line: 2, Column: 5}, End: token.Position{Line: 2, Column: 5}},
					Value: "y",
				},
				Value: &ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "10", Start: token.Position{Line: 2, Column: 9}, End: token.Position{Line: 2, Column: 9}},
				},
			},
		},
	}

	w := ast.CodeWriter{
		Builder: strings.Builder{},
		Mapper:  sourcemap.New(),
	}

	program.WriteTo(&w)

	output := w.String()
	if output != "let x=5;let y=10;" {
		t.Errorf("Output = %q, want %q", output, "let x=5;let y=10;")
	}
}
