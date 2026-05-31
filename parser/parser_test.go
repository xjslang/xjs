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
	LparenToken token.Token
	RparenToken token.Token
	Message     token.Token
}

func (node *MyCustomStmt) Type() string {
	return "MyCustomStmt"
}

func ExampleParser_Init() {
	s := &scanner.Scanner{}
	s.Init([]byte("print('Hello, World!')"))
	p := &parser.Parser{}

	// Declare "middlewares" BEFORE calling Init
	p.UseStmtParser(func(p *parser.Parser, next func() (ast.Node, error)) (_ ast.Node, err error) {
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
			expected: `BinaryExpr
	LeftValue: BinaryExpr
		LeftValue: BasicLit{Value: "1"}
		Operator: "-"
		RightValue: BasicLit{Value: "2"}
	Operator: "-"
	RightValue: BasicLit{Value: "3"}`,
		},
		{
			input: "1 + 2 * (3 + 5) - 4",
			expected: `BinaryExpr
	LeftValue: BinaryExpr
		LeftValue: BasicLit{Value: "1"}
		Operator: "+"
		RightValue: BinaryExpr
			LeftValue: BasicLit{Value: "2"}
			Operator: "*"
			RightValue: ParenExpr
				Value: BinaryExpr
					LeftValue: BasicLit{Value: "3"}
					Operator: "+"
					RightValue: BasicLit{Value: "5"}
	Operator: "-"
	RightValue: BasicLit{Value: "4"}`,
		},
		{
			input: "foo() * 2 + 1",
			expected: `BinaryExpr
	LeftValue: BinaryExpr
		LeftValue: CallExpr
			Function: Ident{Name: "foo"}
		Operator: "*"
		RightValue: BasicLit{Value: "2"}
	Operator: "+"
	RightValue: BasicLit{Value: "1"}`,
		},
		{
			input: "foo(1, 2, 3)",
			expected: `CallExpr
	Function: Ident{Name: "foo"}
	Arguments[0]: BasicLit{Value: "1"}
	Arguments[1]: BasicLit{Value: "2"}
	Arguments[2]: BasicLit{Value: "3"}`,
		},
		{
			input: "2 * (pow(2, 1 + 3) + 4)",
			expected: `BinaryExpr
	LeftValue: BasicLit{Value: "2"}
	Operator: "*"
	RightValue: ParenExpr
		Value: BinaryExpr
			LeftValue: CallExpr
				Function: Ident{Name: "pow"}
				Arguments[0]: BasicLit{Value: "2"}
				Arguments[1]: BinaryExpr
					LeftValue: BasicLit{Value: "1"}
					Operator: "+"
					RightValue: BasicLit{Value: "3"}
			Operator: "+"
			RightValue: BasicLit{Value: "4"}`,
		},
		{
			input: "1 + foo()",
			expected: `BinaryExpr
	LeftValue: BasicLit{Value: "1"}
	Operator: "+"
	RightValue: CallExpr
		Function: Ident{Name: "foo"}`,
		},
		{
			input: "1 + foo()()",
			expected: `BinaryExpr
	LeftValue: BasicLit{Value: "1"}
	Operator: "+"
	RightValue: CallExpr
		Function: CallExpr
			Function: Ident{Name: "foo"}`,
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
