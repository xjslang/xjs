package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type orExpr struct {
	Expr         ast.Node
	FallbackStmt ast.Node
}

func (node *orExpr) Type() string {
	return "orExpr"
}

func TestUseExprParser(t *testing.T) {
	orType := token.RegisterType("or")
	input := "let x = openDb() or exit('failed opening db')"
	s := xjs.NewScanner()
	// the scanner can now scan "or"
	s.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == orType.String() {
			tok.Type = orType
		}
		return tok
	})
	s.Init([]byte(input))
	p := xjs.NewParser()
	// the parser can now parse "or"
	p.UseExprParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		node, err = next()
		if err != nil {
			return
		}
		if p.CurrentToken.Type == orType {
			orNode := &orExpr{Expr: node}
			p.AdvanceToken() // consume "or"
			if orNode.FallbackStmt, err = p.ParseStmt(); err != nil {
				return nil, err
			}
			node = orNode
		}
		return
	})
	p.Init(s)
	result, err := parser.ParseProgram(p)
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.Len(t, result.Stmts, 1)
	require.IsType(t, &js.Let{}, result.Stmts[0])
	stmt := result.Stmts[0].(*js.Let)
	require.IsType(t, &orExpr{}, stmt.Value)
	orVal := stmt.Value.(*orExpr)
	require.IsType(t, &ast.ExprStmt{}, orVal.FallbackStmt)
	fallback := orVal.FallbackStmt.(*ast.ExprStmt)
	require.IsType(t, &ast.CallExpr{}, fallback.Expr)
	expr := fallback.Expr.(*ast.CallExpr)
	require.IsType(t, &ast.Ident{}, expr.Function)
	funcName := expr.Function.(*ast.Ident)
	assert.Equal(t, "exit", funcName.Value.Literal)
}

type notBitwiseExpr struct {
	Operator token.Token
	Value    ast.Node
}

func (node *notBitwiseExpr) Type() string {
	return "notBitwiseExpr"
}

func TestUsePrefixOpParser(t *testing.T) {
	notBitwise := token.RegisterPrefixOp("~")
	input := "1 + ~7"
	s := &scanner.Scanner{}
	s.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		if sc.CurrentChar() == '~' {
			sc.AdvanceChar()
			return token.Token{Type: notBitwise, Literal: "~"}
		}
		return next()
	})
	s.Init([]byte(input))
	p := &parser.Parser{}
	p.UsePrefixOpParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == notBitwise {
			nodeExpr := &notBitwiseExpr{Operator: p.CurrentToken}
			p.AdvanceToken() // consume ~
			if nodeExpr.Value, err = parser.ParseValue(p); err != nil {
				return
			}
			node = nodeExpr
			return
		}
		return next()
	})
	p.Init(s)
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	require.IsType(t, &ast.BinaryExpr{}, result)
	binNode := result.(*ast.BinaryExpr)
	require.IsType(t, &ast.BasicLit{}, binNode.LeftValue)
	require.Equal(t, "1", binNode.LeftValue.(*ast.BasicLit).Value.Literal)
	require.Equal(t, token.PLUS, binNode.Operator.Type)
	require.IsType(t, &notBitwiseExpr{}, binNode.RightValue)
	rightNode := binNode.RightValue.(*notBitwiseExpr)
	require.IsType(t, &ast.BasicLit{}, rightNode.Value)
	require.Equal(t, "7", rightNode.Value.(*ast.BasicLit).Value.Literal)
	require.Equal(t, notBitwise, rightNode.Operator.Type)
}

type powExpr struct {
	LeftValue  ast.Node
	Operator   token.Token
	RightValue ast.Node
}

func (node *powExpr) Type() string {
	return "powExpr"
}

func TestUseInfixOpParser(t *testing.T) {
	powType := token.RegisterInfixOp("^", token.MULTIPLY.Precedence()+1)
	input := "1+5^2"
	s := &scanner.Scanner{}
	s.UseScanner(func(s *scanner.Scanner, next func() token.Token) token.Token {
		if s.CurrentChar() == '^' {
			s.AdvanceChar()
			return token.Token{Type: powType, Literal: "^"}
		}
		return next()
	})
	s.Init([]byte(input))
	p := &parser.Parser{}
	p.UseInfixOpParser(func(p *parser.Parser, leftVal ast.Node, next func(ast.Node) (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == powType {
			powNode := &powExpr{LeftValue: leftVal, Operator: p.CurrentToken}
			p.AdvanceToken() // consume ^
			if powNode.RightValue, err = parser.ParseRightValue(p, powType.Precedence()); err != nil {
				return
			}
			node = powNode
			return
		}
		return next(leftVal)
	})
	p.Init(s)
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.IsType(t, &ast.BinaryExpr{}, result)
	binNode := result.(*ast.BinaryExpr)
	require.IsType(t, &ast.BasicLit{}, binNode.LeftValue)
	require.Equal(t, "1", binNode.LeftValue.(*ast.BasicLit).Value.Literal)
	require.Equal(t, token.PLUS, binNode.Operator.Type)
	require.IsType(t, &powExpr{}, binNode.RightValue)
	powNode := binNode.RightValue.(*powExpr)
	require.IsType(t, &ast.BasicLit{}, powNode.LeftValue)
	require.Equal(t, "5", powNode.LeftValue.(*ast.BasicLit).Value.Literal)
	require.Equal(t, powType, powNode.Operator.Type)
	require.IsType(t, &ast.BasicLit{}, powNode.RightValue)
	require.Equal(t, "2", powNode.RightValue.(*ast.BasicLit).Value.Literal)
}

type factorialExpr struct {
	Operator token.Token
	Value    ast.Node
}

func (node *factorialExpr) Type() string {
	return "factorialExpr"
}

func TestUseInfixOpParser_postfix(t *testing.T) {
	facTyp := token.RegisterInfixOp("!", -1)
	input := "5! + 1"
	s := &scanner.Scanner{}
	s.UseScanner(func(s *scanner.Scanner, next func() token.Token) token.Token {
		if s.CurrentChar() == '!' {
			s.AdvanceChar()
			return token.Token{Type: facTyp, Literal: "!"}
		}
		return next()
	})
	s.Init([]byte(input))
	p := &parser.Parser{}
	p.UseInfixOpParser(func(p *parser.Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == facTyp {
			leftVal = &factorialExpr{Operator: p.CurrentToken, Value: leftVal}
			p.AdvanceToken() // consume !
		}
		return next(leftVal)
	})
	p.Init(s)
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.IsType(t, &ast.BinaryExpr{}, result)
	binNode := result.(*ast.BinaryExpr)
	require.IsType(t, &factorialExpr{}, binNode.LeftValue)
	leftNode := binNode.LeftValue.(*factorialExpr)
	require.IsType(t, &ast.BasicLit{}, leftNode.Value)
	leftVal := leftNode.Value.(*ast.BasicLit)
	require.Equal(t, "5", leftVal.Value.Literal)
	require.Equal(t, facTyp, leftNode.Operator.Type)
	require.IsType(t, &ast.BasicLit{}, binNode.RightValue)
	rightVal := binNode.RightValue.(*ast.BasicLit)
	require.Equal(t, "1", rightVal.Value.Literal)
}
