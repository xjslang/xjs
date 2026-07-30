package parser_test

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
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

func ExampleBuilder_Build() {
	s := scanner.NewBuilder().Build([]byte("print('Hello, World!')"))
	p := parser.NewBuilder().
		UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (_ ast.Stmt, err error) {
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
		}).
		Build(s)

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

func TestLookahead(t *testing.T) {
	parser1 := func(p *parser.Parser) (node *js.Variable, err error) {
		node = &js.Variable{}
		if node.Token, err = p.ExpectLiteral("a"); err != nil {
			return
		}
		return
	}

	parser2 := func(p *parser.Parser) (node *js.Variable, err error) {
		node = &js.Variable{}
		if node.Token, err = p.ExpectLiteral("b"); err != nil {
			return
		}
		return
	}

	t.Run("resolved", func(t *testing.T) {
		input := "b"
		sc := scanner.NewBuilder().Build([]byte(input))
		p := parser.NewBuilder().Build(sc)
		result, err := parser.Switch(p, func(p *parser.Parser) (*js.Variable, error) {
			return parser1(p)
		}, func(p *parser.Parser) (*js.Variable, error) {
			return parser2(p)
		})
		require.NoError(t, err)
		testutil.AssertTokens(t, []token.Token{result.Token}, []token.Token{
			{Type: token.IDENT, Literal: "b"},
		})
		require.Equal(t, token.EOF, p.CurrentToken.Type)
	})

	t.Run("not resolved", func(t *testing.T) {
		input := "c"
		sc := scanner.NewBuilder().Build([]byte(input))
		p := parser.NewBuilder().Build(sc)
		_, err := parser.Switch(p, func(p *parser.Parser) (*js.Variable, error) {
			return parser1(p)
		}, func(p *parser.Parser) (*js.Variable, error) {
			return parser2(p)
		})
		require.Error(t, err)
		require.Equal(t, token.IDENT, p.CurrentToken.Type)
		require.Equal(t, "c", p.CurrentToken.Literal)
	})
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
			p := xjs.PluginBuilder().Build([]byte(test.input))
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
			p := xjs.PluginBuilder().Build([]byte(test.input))
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
		p := xjs.PluginBuilder().Build([]byte(input))
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
					{Type: token.LINE_COMMENT, Literal: "// comment before {\n"},
					{Type: token.NEWLINE, Literal: "\n"},
				}},
				{Type: token.RBRACE, Literal: "}", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: "// comment before }\n"},
					{Type: token.BLOCK_COMMENT, Literal: "/* block comment */"},
				}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
	t.Run("grouped expression", func(t *testing.T) {
		input := `// comment before
	(1 + 2// comment after
	)`
		p := xjs.PluginBuilder().Build([]byte(input))
		result, err := js.ParseGroupExpr(p)
		if err != nil {
			t.Fatal(err)
		}
		testutil.AssertTokens(
			t,
			[]token.Token{result.Layout.Lparen, result.Layout.Rparen},
			[]token.Token{
				{Type: token.LPAREN, Literal: "(", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: "// comment before\n"},
				}},
				{Type: token.RPAREN, Literal: ")", LeadingTrivia: []token.Token{
					{Type: token.LINE_COMMENT, Literal: "// comment after\n"},
				}},
			},
			testutil.CompareLeadingTrivia(),
		)
	})
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
				p := xjs.PluginBuilder().Build([]byte(input))
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

func TestExpectWith(t *testing.T) {
	scanRegex := func(sc token.Scanner) (string, error) {
		sb := strings.Builder{}
		if sc.CurrentChar() != '/' {
			return "", errors.New("/ expected")
		}
		sb.WriteRune(sc.CurrentChar())
		sc.AdvanceChar() // consume /
		for {
			if sc.CurrentChar() == '\\' {
				sb.WriteRune(sc.CurrentChar())
				sc.AdvanceChar()
				if sc.CurrentChar() == '/' {
					sb.WriteRune(sc.CurrentChar())
					sc.AdvanceChar()
					continue
				}
			}
			if sc.CurrentChar() == '/' {
				sb.WriteRune(sc.CurrentChar())
				sc.AdvanceChar()
				break
			} else if sc.CurrentChar() == scanner.EOF || sc.CurrentChar() == '\n' || sc.CurrentChar() == '\r' {
				return sb.String(), errors.New("unexpected end of line")
			}
			sb.WriteRune(sc.CurrentChar())
			sc.AdvanceChar()
		}
		flags := []rune("dgimsuvy")
		for slices.Contains(flags, sc.CurrentChar()) {
			sb.WriteRune(sc.CurrentChar())
			sc.AdvanceChar()
		}
		if c := sc.CurrentChar(); c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != scanner.EOF {
			return "", errors.New("unknown flag " + string(c))
		}
		return sb.String(), nil
	}

	input := `// comment
	/lorem ipsum dolor/gd
	'lorem ipsum'`
	sc := scanner.NewBuilder().Build([]byte(input))
	p := parser.NewBuilder().Build(sc)
	tok, err := p.ExpectWith(scanRegex)
	require.NoError(t, err)
	testutil.AssertTokens(t, []token.Token{
		tok,
		p.CurrentToken,
		p.PeekToken,
	}, []token.Token{
		{Type: token.DIVIDE, Literal: "/lorem ipsum dolor/gd"},
		{Type: token.DIVIDE, Literal: "/lorem ipsum dolor/gd"},
		{Type: token.STRING, Literal: "'lorem ipsum'"},
	})
}
