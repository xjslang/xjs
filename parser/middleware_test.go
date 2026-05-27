package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs/ast"
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
	s := &scanner.Scanner{}
	// the scanner can now scan "or"
	s.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == orType.String() {
			tok.Type = orType
		}
		return tok
	})
	s.Init([]byte(input))
	p := &parser.Parser{}
	// the parser can now parse "or"
	p.UseExprParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		node, err = next()
		if err != nil {
			return
		}
		if p.CurrentToken.Type == orType {
			orNode := &orExpr{Expr: node}
			p.AdvanceToken() // consume "or"
			orNode.FallbackStmt, err = p.ParseStmt()
			if err != nil {
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
	require.IsType(t, &ast.LetStmt{}, result.Stmts[0])
	stmt := result.Stmts[0].(*ast.LetStmt)
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

type factorialExpr struct {
	Operator token.Token
	Value    ast.Node
}

func (node *factorialExpr) Type() string {
	return "factorialExpr"
}

func TestUseUnaryExprParser(t *testing.T) {
	facType := token.RegisterUnaryOperator("¡")
	input := "1 + ¡7"
	s := &scanner.Scanner{}
	s.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		if sc.CurrentChar == '¡' {
			sc.AdvanceChar()
			return token.Token{Type: facType, Literal: "¡"}
		}
		return next()
	})
	s.Init([]byte(input))
	p := &parser.Parser{}
	p.UseUnaryExprParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == facType {
			factorialNode := &factorialExpr{Operator: p.CurrentToken}
			p.AdvanceToken()
			if factorialNode.Value, err = parser.ParseValue(p); err != nil {
				return
			}
			node = factorialNode
			return
		}
		return next()
	})
	p.Init(s)
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	require.IsType(t, &ast.InfixExpr{}, result)
	infixNode := result.(*ast.InfixExpr)
	require.IsType(t, &ast.BasicLit{}, infixNode.LeftValue)
	require.Equal(t, "1", infixNode.LeftValue.(*ast.BasicLit).Value.Literal)
	require.Equal(t, token.PLUS, infixNode.Operator.Type)
	require.IsType(t, &factorialExpr{}, infixNode.RightValue)
	facNode := infixNode.RightValue.(*factorialExpr)
	require.IsType(t, &ast.BasicLit{}, facNode.Value)
	require.Equal(t, "7", facNode.Value.(*ast.BasicLit).Value.Literal)
	require.Equal(t, facType, facNode.Operator.Type)
}

type powExpr struct {
	LeftValue  ast.Node
	Operator   token.Token
	RightValue ast.Node
}

func (node *powExpr) Type() string {
	return "powExpr"
}

func TestUseBinExprParser(t *testing.T) {
	powType := token.RegisterBinaryOperator("^", token.MULTIPLY.Precedence()+1)
	input := "1+5^2"
	s := &scanner.Scanner{}
	s.UseScanner(func(s *scanner.Scanner, next func() token.Token) token.Token {
		if s.CurrentChar == '^' {
			s.AdvanceChar()
			return token.Token{Type: powType, Literal: "^"}
		}
		return next()
	})
	s.Init([]byte(input))
	p := &parser.Parser{}
	p.UseBinExprParser(func(p *parser.Parser, leftVal ast.Node, next func(ast.Node) (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == powType {
			powNode := &powExpr{LeftValue: leftVal, Operator: p.CurrentToken}
			if powNode.RightValue, err = parser.ParseRightExpr(p); err != nil {
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
	require.IsType(t, &ast.InfixExpr{}, result)
	infixNode := result.(*ast.InfixExpr)
	require.IsType(t, &ast.BasicLit{}, infixNode.LeftValue)
	require.Equal(t, "1", infixNode.LeftValue.(*ast.BasicLit).Value.Literal)
	require.Equal(t, token.PLUS, infixNode.Operator.Type)
	require.IsType(t, &powExpr{}, infixNode.RightValue)
	powNode := infixNode.RightValue.(*powExpr)
	require.IsType(t, &ast.BasicLit{}, powNode.LeftValue)
	require.Equal(t, "5", powNode.LeftValue.(*ast.BasicLit).Value.Literal)
	require.Equal(t, powType, powNode.Operator.Type)
	require.IsType(t, &ast.BasicLit{}, powNode.RightValue)
	require.Equal(t, "2", powNode.RightValue.(*ast.BasicLit).Value.Literal)
}
