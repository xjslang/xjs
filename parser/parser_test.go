package parser_test

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var updateGoldenFiles bool

type MyCustomStmt struct {
	ast.BaseStmt
	LparenToken token.Token
	RparenToken token.Token
	Message     token.Token
}

func ExampleParser_Init() {
	s := &scanner.Scanner{}
	s.Init([]byte("print('Hello, World!')"))
	p := &parser.Parser{}

	// Declare "middlewares" BEFORE calling Init
	p.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (_ ast.Stmt, err error) {
		if p.CurrentToken.Type == token.IDENT && p.CurrentToken.Literal == "print" {
			p.AdvanceToken()
			node := &MyCustomStmt{}
			if node.LparenToken, err = p.Expect(token.LPAREN); err != nil { // expect (
				return
			}
			if node.Message, err = p.Expect(token.STRING); err != nil { // expect a string
				return
			}
			if node.RparenToken, err = p.Expect(token.RPAREN); err != nil { // expect )
				return
			}
			return node, nil
		}
		return next() // Delegate to the "next" middleware
	})
	p.Init(s)

	// Now you can use the parser
	result, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}
	stmt := result.Stmts[0].(*MyCustomStmt)
	fmt.Println(stmt.Message.Literal)
	// Output: 'Hello, World!'
}

func TestMain(m *testing.M) {
	flag.BoolVar(&updateGoldenFiles, "update", false, "update golden files")
	flag.Parse()
	os.Exit(m.Run())
}

func TestExprs(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: "1 - 2 - 3",
			expected: `*js.BinaryExpr
	Left: *js.BinaryExpr
		Left: *js.Literal{Value: "1"}
		Op: "-"
		Right: *js.Literal{Value: "2"}
	Op: "-"
	Right: *js.Literal{Value: "3"}`,
		},
		{
			input: "1 + 2 * (3 + 5) - 4",
			expected: `*js.BinaryExpr
	Left: *js.BinaryExpr
		Left: *js.Literal{Value: "1"}
		Op: "+"
		Right: *js.BinaryExpr
			Left: *js.Literal{Value: "2"}
			Op: "*"
			Right: *js.GroupExpr
				Value: *js.BinaryExpr
					Left: *js.Literal{Value: "3"}
					Op: "+"
					Right: *js.Literal{Value: "5"}
	Op: "-"
	Right: *js.Literal{Value: "4"}`,
		},
		{
			input: "foo() * 2 + 1",
			expected: `*js.BinaryExpr
	Left: *js.BinaryExpr
		Left: *js.CallExpr
			Callee: *js.Variable{Name: "foo"}
		Op: "*"
		Right: *js.Literal{Value: "2"}
	Op: "+"
	Right: *js.Literal{Value: "1"}`,
		},
		{
			input: "foo(1, 2, 3)",
			expected: `*js.CallExpr
	Callee: *js.Variable{Name: "foo"}
	Args[0]: *js.Literal{Value: "1"}
	Args[1]: *js.Literal{Value: "2"}
	Args[2]: *js.Literal{Value: "3"}`,
		},
		{
			input: "2 * (pow(2, 1 + 3) + 4)",
			expected: `*js.BinaryExpr
	Left: *js.Literal{Value: "2"}
	Op: "*"
	Right: *js.GroupExpr
		Value: *js.BinaryExpr
			Left: *js.CallExpr
				Callee: *js.Variable{Name: "pow"}
				Args[0]: *js.Literal{Value: "2"}
				Args[1]: *js.BinaryExpr
					Left: *js.Literal{Value: "1"}
					Op: "+"
					Right: *js.Literal{Value: "3"}
			Op: "+"
			Right: *js.Literal{Value: "4"}`,
		},
		{
			input: "1 + foo()",
			expected: `*js.BinaryExpr
	Left: *js.Literal{Value: "1"}
	Op: "+"
	Right: *js.CallExpr
		Callee: *js.Variable{Name: "foo"}`,
		},
		{
			input: "1 + foo()()",
			expected: `*js.BinaryExpr
	Left: *js.Literal{Value: "1"}
	Op: "+"
	Right: *js.CallExpr
		Callee: *js.CallExpr
			Callee: *js.Variable{Name: "foo"}`,
		},
	}
	for i, test := range tests {
		t.Run("exp "+strconv.Itoa(i), func(t *testing.T) {
			p := xjs.NewBuilder().Build([]byte(test.input))
			result, err := js.ParseExpr(p)
			if err != nil {
				t.Fatal(err)
			}
			if got := testutil.NodeString(result); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}

func TestMalformedExpr(t *testing.T) {
	t.Run("block", func(t *testing.T) {
		tests := []struct {
			input       string
			expectedErr string
		}{
			{"let x = 100 }", "{ expected"},
			{"{ let x = 100", "} expected"},
		}
		for i, test := range tests {
			p := xjs.NewBuilder().Build([]byte(test.input))
			_, err := js.ParseBlockStmt(p)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			if got := err.Error(); !strings.HasSuffix(got, test.expectedErr) {
				t.Fatalf("%d: Expected %q, got %q", i, test.expectedErr, got)
			}
		}
	})
	t.Run("grouped expression", func(t *testing.T) {
		tests := []struct {
			input       string
			expectedErr string
		}{
			{"1 + 2)", "( expected"},
			{"(1 + 2", ") expected"},
		}
		for i, test := range tests {
			p := xjs.NewBuilder().Build([]byte(test.input))
			_, err := js.ParseGroupExpr(p)
			if err == nil {
				t.Fatal("Expected an error, got nil")
			}
			if got := err.Error(); !strings.HasSuffix(got, test.expectedErr) {
				t.Fatalf("%d: Expected error to be %q, got %q", i, test.expectedErr, got)
			}
		}
	})
}

func TestKeysAreSaved(t *testing.T) {
	t.Run("block", func(t *testing.T) {
		input := `
		// comment before {

		{
		let x = 100
		let y = 200 // comment before }
		/* block comment */ }`
		p := xjs.NewBuilder().Build([]byte(input))
		result, err := js.ParseBlockStmt(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.Layout.Lbrace, result.Layout.Rbrace},
			[]token.Token{
				{Type: token.LBRACE, Literal: "{", LeadingTrivia: []token.Token{
					{Type: token.NEWLINE, Literal: "\n"},
					{Type: token.LINE_COMMENT, Literal: " comment before {\n"},
					{Type: token.NEWLINE, Literal: "\n"},
				}},
				{Type: token.RBRACE, Literal: "}", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment before }\n"},
					{Type: token.BLOCK_COMMENT, Literal: " block comment "},
				}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
	t.Run("grouped expression", func(t *testing.T) {
		input := `// comment before
	(1 + 2// comment after
	)`
		p := xjs.NewBuilder().Build([]byte(input))
		result, err := js.ParseGroupExpr(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.Layout.Lparen, result.Layout.Rparen},
			[]token.Token{
				{Type: token.LPAREN, Literal: "(", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment before\n"},
				}},
				{Type: token.RPAREN, Literal: ")", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: " comment after\n"},
				}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
}

func TestStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: "log()",
			expected: `*js.ExprStmt
	Expr: *js.CallExpr
		Callee: *js.Variable{Name: "log"}`,
		},
		{
			input: "log(1)",
			expected: `*js.ExprStmt
	Expr: *js.CallExpr
		Callee: *js.Variable{Name: "log"}
		Args[0]: *js.Literal{Value: "1"}`,
		},
		{
			input: "log(1, 2)",
			expected: `*js.ExprStmt
	Expr: *js.CallExpr
		Callee: *js.Variable{Name: "log"}
		Args[0]: *js.Literal{Value: "1"}
		Args[1]: *js.Literal{Value: "2"}`,
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			p := xjs.NewBuilder().Build([]byte(test.input))
			node, err := js.ParseExprStmt(p)
			if err != nil {
				t.Fatal(err)
			}
			if got := testutil.NodeString(node); got != test.expected {
				t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", test.expected, got)
			}
		})
	}
}

func TestInvalidTokenAfterNewline(t *testing.T) {
	tests := []string{"\n%", "let\n%", "let x\n%", "let y =\n%", "let x =\nlet y = 1"}
	for i := range 2 {
		for j, test := range tests {
			t.Run(fmt.Sprintf("test %d%d", i, j), func(t *testing.T) {
				var input string
				if i > 0 {
					input = fmt.Sprintf("{%s}", test)
				} else {
					input = test
				}
				p := xjs.NewBuilder().Build([]byte(input))
				var err error
				if i > 0 {
					_, err = js.ParseBlockStmt(p)
				} else {
					_, err = js.ParseProgram(p)
				}
				var errList parser.ErrorList
				require.ErrorAs(t, err, &errList)
				require.NotEmpty(t, errList)
			})
		}
	}
}
