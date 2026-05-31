package testutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/token"
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

func AssertTokens(t *testing.T, toks, expectedToks []token.Token, opts ...TokenCompareOption) {
	t.Helper()
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
				for j, expectedTrivia := range expectedTok.LeadingTrivia {
					trivia := tok.LeadingTrivia[j]
					if trivia.Type != expectedTrivia.Type {
						t.Errorf("token %d: expected trivia type to be %v, got %v", i, expectedTrivia.Type, trivia.Type)
					} else if trivia.Literal != expectedTrivia.Literal {
						t.Errorf("token %d: expected trivia to be %q, got %q", i, expectedTrivia.Literal, trivia.Literal)
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
		fmt.Fprint(s, node.Type())
		switch v := node.(type) {
		case *js.Ident:
			fmt.Fprintf(s, "{Name: %q}", v.Name.Literal)
		case *js.Program:
			for _, stmt := range v.Stmts {
				fmt.Fprintf(s, "\n%s%s", indent, print(stmt))
			}
		case *js.BlockStmt:
			for _, stmt := range v.Stmts {
				fmt.Fprintf(s, "\n%s%s", indent, print(stmt))
			}
		case *js.Stmt:
			fmt.Fprintf(s, "\n%sExpr: %s", indent, print(v.Expr))
		case *js.LetStmt:
			fmt.Fprintf(s, "\n%sName: %s", indent, v.Name.Literal)
			fmt.Fprintf(s, "\n%sValue: %s", indent, print(v.Value))
		case *js.FunctionDecl:
			fmt.Fprintf(s, "\n%sName: %s", indent, v.Name.Literal)
			fmt.Fprintf(s, "\n%sBody: %s", indent, print(v.Body))
		case *js.ParenExpr:
			fmt.Fprintf(s, "\n%sValue: %s", indent, print(v.Value))
		case *js.CallExpr:
			fmt.Fprintf(s, "\n%sFunction: %s", indent, print(v.Function))
			for i, arg := range v.Arguments {
				fmt.Fprintf(s, "\n%sArguments[%d]: %s", indent, i, print(arg))
			}
		case *js.BinaryExpr:
			fmt.Fprintf(s, "\n%sLeftValue: %s", indent, print(v.LeftValue))
			fmt.Fprintf(s, "\n%sOperator: %q", indent, v.Operator.Type.String())
			fmt.Fprintf(s, "\n%sRightValue: %s", indent, print(v.RightValue))
		case *js.Literal:
			fmt.Fprintf(s, "{Value: %q}", v.Value.Literal)
		}
		return s.String()
	}
	return print(node)
}
