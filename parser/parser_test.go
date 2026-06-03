package parser_test

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"

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
			Right: *js.ParenExpr
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
	Right: *js.ParenExpr
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
