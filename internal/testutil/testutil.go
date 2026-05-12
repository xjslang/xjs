package testutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
)

type tokenCompareConfig struct {
	afterNewline  bool
	leadingTrivia bool
	tokenPosition bool
}

type TokenCompareOption func(cfg *tokenCompareConfig)

func CompareAfterNewline() TokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.afterNewline = true
	}
}

func CompareLeadingTrivia() TokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.leadingTrivia = true
	}
}

func CompareTokenPosition() TokenCompareOption {
	return func(cfg *tokenCompareConfig) {
		cfg.tokenPosition = true
	}
}

func AssertTokens(t *testing.T, toks []scanner.Token, expectedToks []scanner.Token, opts ...TokenCompareOption) {
	cfg := &tokenCompareConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if len(toks) != len(expectedToks) {
		t.Fatalf("Expect len(toks) = %d, got %d", len(toks), len(expectedToks))
	}
	for i, expectedTok := range expectedToks {
		tok := toks[i]
		switch {
		case tok.Type != expectedTok.Type:
			t.Errorf("token %d: expected type %v, got %v", i, expectedTok.Type, tok.Type)
		case tok.Literal != expectedTok.Literal:
			t.Errorf("token %d: expected %q, got %q", i, expectedTok.Literal, tok.Literal)
		case cfg.afterNewline && tok.AfterNewline != expectedTok.AfterNewline:
			t.Errorf("token %d: expected AfterNewline to be %t, got %t", i, expectedTok.AfterNewline, tok.AfterNewline)
		case cfg.leadingTrivia:
			if len(tok.LeadingTrivia) != len(expectedTok.LeadingTrivia) {
				t.Errorf("token %d: expected %d leading trivia lines, got %d", i, len(expectedTok.LeadingTrivia), len(tok.LeadingTrivia))
			} else {
				for j, line := range expectedTok.LeadingTrivia {
					if tok.LeadingTrivia[j] != line {
						t.Errorf("token %d: expected %q leading trivia line, got %q", i, line, tok.LeadingTrivia[j])
					}
				}
			}
		case cfg.tokenPosition && (tok.Line != expectedTok.Line || tok.Column != expectedTok.Column):
			t.Errorf("token %d: expected position to be (%d, %d), got (%d, %d)", i, expectedTok.Line, expectedTok.Column, tok.Line, tok.Column)
		}
	}
}

func NodeString(node ast.Node) string {
	indentLevel := 0
	var print func(node ast.Node) string
	print = func(node ast.Node) string {
		s := &strings.Builder{}
		indentLevel++
		defer func() {
			indentLevel--
		}()
		indent := strings.Repeat("\t", indentLevel)
		fmt.Print(node.Type())
		switch v := node.(type) {
		case *ast.Block:
			for _, stmt := range v.Statements {
				fmt.Fprintf(s, "\n%s%s", indent, print(stmt))
			}
		case *ast.Let:
			fmt.Fprintf(s, "\n%sName: %s", indent, v.Name.Literal)
			fmt.Fprintf(s, "\n%sValue: %s", indent, print(v.Value))
		case *ast.Function:
			fmt.Fprintf(s, "\n%sName: %s", indent, v.Name.Literal)
			fmt.Fprintf(s, "\n%sBody: %s", indent, print(v.Body))
		case *ast.GroupedExpression:
			fmt.Fprintf(s, "\n%sValue: %s", indent, print(v.Value))
		case *ast.InfixOperator:
			fmt.Fprintf(s, "\n%sLeftValue: %s", indent, print(v.LeftValue))
			fmt.Fprintf(s, "\n%sOperator: %q", indent, v.Operator.Type.String())
			fmt.Fprintf(s, "\n%sRightValue: %s", indent, print(v.RightValue))
		case *ast.Integer:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		case *ast.String:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		case *ast.Boolean:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		case *ast.Ident:
			fmt.Fprintf(s, "{Value: %q}", v.Value)
		}
		return s.String()
	}
	return print(node)
}

func Parse(input string) (*ast.Block, error) {
	sc := &scanner.Scanner{}
	sc.Init([]byte(input))
	p := &parser.Parser{}
	p.Init(sc)
	return parser.ParseProgram(p)
}
