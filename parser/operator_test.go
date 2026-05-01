package parser

import (
	"testing"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/internal/testutil"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type PowExpression struct {
	LeftValue  ast.Expression
	Operator   token.Token
	RightValue ast.Expression
}

func (node *PowExpression) PrintTo(p *printer.Printer) {
	p.PrintString("Math.pow(")
	node.LeftValue.PrintTo(p)
	p.PrintString(", ")
	node.RightValue.PrintTo(p)
	p.PrintRune(')')
}

func TestRegisterInfixOperator(t *testing.T) {
	l := &lexer.Lexer{}
	powType := token.RegisterType("**")
	l.UseTokenizer(func(l *lexer.Lexer, next func() token.Token) token.Token {
		if l.CurrentChar == '*' && l.PeekChar() == '*' {
			l.AdvanceChar()
			l.AdvanceChar()
			return token.Token{Type: powType, Literal: powType.String()}
		}
		return next()
	})
	l.Init([]byte("1 + 2 ** 3"))

	p := Parser{}
	err := p.RegisterInfixOperator(powType, 3, func(op token.Token, left, right ast.Expression) ast.Expression {
		return &PowExpression{LeftValue: left, Operator: op, RightValue: right}
	})
	if err != nil {
		t.Fatal(err)
	}
	p.Init(l)
	exp, err := p.ParseExpression()
	if err != nil {
		t.Fatal(err)
	}

	expected := `*ast.InfixOperator
	LeftValue: *ast.IntegerLiteral{Value: "1"}
	Operator: "+"
	RightValue: *parser.PowExpression`
	if result := testutil.NodeString(exp); result != expected {
		t.Errorf("Invalid node:\nExpected:\n%s\nGot:\n%s", expected, result)
	}

	prt := printer.New()
	exp.PrintTo(prt)
	expected = `1 + Math.pow(2, 3)`
	if result := prt.String(); result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	errorTests := []struct {
		tok        token.TokenType
		precedence int
		fn         func(op token.Token, left, right ast.Expression) ast.Expression
		expected   string
	}{
		{powType, -1, nil, "negative precedence"},
		{powType, 1, nil, "nil function"},
	}
	for _, test := range errorTests {
		t.Run(test.expected, func(t *testing.T) {
			p := Parser{}
			err := p.RegisterInfixOperator(test.tok, test.precedence, test.fn)
			if err == nil {
				t.Error("Expected an error, got nil")
			} else if result := err.Error(); result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}

	t.Run("operator already registered", func(t *testing.T) {
		p := Parser{}
		for i := range 2 {
			err := p.RegisterInfixOperator(powType, 3, func(op token.Token, left, right ast.Expression) ast.Expression {
				return &PowExpression{LeftValue: left, Operator: op, RightValue: right}
			})
			if i == 0 && err != nil {
				t.Error(err)
			}
			if i > 0 {
				expected := "operator already registered"
				if err == nil {
					t.Errorf("Expected an error")
				} else if result := err.Error(); result != expected {
					t.Errorf("Expected %q, got %q", expected, result)
				}
			}
		}
	})
}
