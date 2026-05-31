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
	b := xjs.NewBuilder()
	// the scanner can now scan "or"
	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == orType.String() {
			tok.Type = orType
		}
		return tok
	})
	// the parser can now parse "or"
	b.UseExprParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		node, err = next()
		if err != nil {
			return
		}
		if p.CurrentToken.Type == orType {
			orNode := &orExpr{Expr: node}
			p.AdvanceToken() // consume "or"
			if orNode.FallbackStmt, err = p.ParseExprStmt(); err != nil {
				return nil, err
			}
			node = orNode
		}
		return
	})
	p := b.Build([]byte(input))
	result, err := js.ParseProgram(p)
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.Len(t, result.Stmts, 1)
	require.IsType(t, &js.LetStmt{}, result.Stmts[0])
	stmt := result.Stmts[0].(*js.LetStmt)
	require.IsType(t, &orExpr{}, stmt.Value)
	orVal := stmt.Value.(*orExpr)
	require.IsType(t, &js.ExprStmt{}, orVal.FallbackStmt)
	fallback := orVal.FallbackStmt.(*js.ExprStmt)
	require.IsType(t, &js.CallExpr{}, fallback.Expr)
	expr := fallback.Expr.(*js.CallExpr)
	require.IsType(t, &js.Ident{}, expr.Function)
	funcName := expr.Function.(*js.Ident)
	assert.Equal(t, "exit", funcName.Name.Literal)
}

type notBitwiseExpr struct {
	Operator token.Token
	Value    ast.Node
}

func (node *notBitwiseExpr) Type() string {
	return "notBitwiseExpr"
}

func TestUseUnaryParser(t *testing.T) {
	notBitwise := token.RegisterUnaryOp("~")
	input := "1 + ~7"
	b := xjs.NewBuilder()
	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		if sc.CurrentChar() == '~' {
			sc.AdvanceChar()
			return token.Token{Type: notBitwise, Literal: "~"}
		}
		return next()
	})
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == notBitwise {
			nodeExpr := &notBitwiseExpr{Operator: p.CurrentToken}
			p.AdvanceToken() // consume ~
			if nodeExpr.Value, err = js.ParseValue(p); err != nil {
				return
			}
			node = nodeExpr
			return
		}
		return next()
	})
	p := b.Build([]byte(input))
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	require.IsType(t, &js.BinaryExpr{}, result)
	binNode := result.(*js.BinaryExpr)
	require.IsType(t, &js.Literal{}, binNode.LeftValue)
	require.Equal(t, "1", binNode.LeftValue.(*js.Literal).Value.Literal)
	require.Equal(t, token.PLUS, binNode.Operator.Type)
	require.IsType(t, &notBitwiseExpr{}, binNode.RightValue)
	rightNode := binNode.RightValue.(*notBitwiseExpr)
	require.IsType(t, &js.Literal{}, rightNode.Value)
	require.Equal(t, "7", rightNode.Value.(*js.Literal).Value.Literal)
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

func TestUseBinaryParser(t *testing.T) {
	powType := token.RegisterBinaryOp("^", token.MULTIPLY.Precedence()+1)
	input := "1+5^2"
	b := xjs.NewBuilder()
	b.UseScanner(func(s *scanner.Scanner, next func() token.Token) token.Token {
		if s.CurrentChar() == '^' {
			s.AdvanceChar()
			return token.Token{Type: powType, Literal: "^"}
		}
		return next()
	})
	b.UseBinaryParser(func(p *parser.Parser, leftVal ast.Node, next func(ast.Node) (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == powType {
			powNode := &powExpr{LeftValue: leftVal, Operator: p.CurrentToken}
			p.AdvanceToken() // consume ^
			if powNode.RightValue, err = js.ParseRightExpr(p, powType.Precedence()); err != nil {
				return
			}
			node = powNode
			return
		}
		return next(leftVal)
	})
	p := b.Build([]byte(input))
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.IsType(t, &js.BinaryExpr{}, result)
	binNode := result.(*js.BinaryExpr)
	require.IsType(t, &js.Literal{}, binNode.LeftValue)
	require.Equal(t, "1", binNode.LeftValue.(*js.Literal).Value.Literal)
	require.Equal(t, token.PLUS, binNode.Operator.Type)
	require.IsType(t, &powExpr{}, binNode.RightValue)
	powNode := binNode.RightValue.(*powExpr)
	require.IsType(t, &js.Literal{}, powNode.LeftValue)
	require.Equal(t, "5", powNode.LeftValue.(*js.Literal).Value.Literal)
	require.Equal(t, powType, powNode.Operator.Type)
	require.IsType(t, &js.Literal{}, powNode.RightValue)
	require.Equal(t, "2", powNode.RightValue.(*js.Literal).Value.Literal)
}

type factorialExpr struct {
	Operator token.Token
	Value    ast.Node
}

func (node *factorialExpr) Type() string {
	return "factorialExpr"
}

func TestUseBinaryParser_postfix(t *testing.T) {
	facTyp := token.RegisterBinaryOp("!", -1)
	input := "5! + 1"
	b := xjs.NewBuilder()
	b.UseScanner(func(s *scanner.Scanner, next func() token.Token) token.Token {
		if s.CurrentChar() == '!' {
			s.AdvanceChar()
			return token.Token{Type: facTyp, Literal: "!"}
		}
		return next()
	})
	b.UseBinaryParser(func(p *parser.Parser, leftVal ast.Node, next func(leftVal ast.Node) (ast.Node, error)) (node ast.Node, err error) {
		if p.CurrentToken.Type == facTyp {
			leftVal = &factorialExpr{Operator: p.CurrentToken, Value: leftVal}
			p.AdvanceToken() // consume !
		}
		return next(leftVal)
	})
	p := b.Build([]byte(input))
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.IsType(t, &js.BinaryExpr{}, result)
	binNode := result.(*js.BinaryExpr)
	require.IsType(t, &factorialExpr{}, binNode.LeftValue)
	leftNode := binNode.LeftValue.(*factorialExpr)
	require.IsType(t, &js.Literal{}, leftNode.Value)
	leftVal := leftNode.Value.(*js.Literal)
	require.Equal(t, "5", leftVal.Value.Literal)
	require.Equal(t, facTyp, leftNode.Operator.Type)
	require.IsType(t, &js.Literal{}, binNode.RightValue)
	rightVal := binNode.RightValue.(*js.Literal)
	require.Equal(t, "1", rightVal.Value.Literal)
}
