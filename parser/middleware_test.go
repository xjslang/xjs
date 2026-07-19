package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/jsextended"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type orExpr struct {
	ast.BaseExpr
	Value        ast.Expr
	FallbackStmt ast.Stmt
}

func TestUseExprParser(t *testing.T) {
	orType := token.RegisterType("or")
	input := "let x = openDb() or exit('failed opening db')"
	b := xjs.PluginBuilder().Install(jsextended.Plugin)
	// the scanner can now scan "or"
	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		if tok.Type == token.IDENT && tok.Literal == orType.String() {
			tok.Type = orType
		}
		return
	})
	// the parser can now parse "or"
	b.UseExprParser(func(p *parser.Parser, next func() (ast.Expr, error)) (node ast.Expr, err error) {
		node, err = next()
		if err != nil {
			return
		}
		if p.CurrentToken.Type == orType {
			orNode := &orExpr{Value: node}
			p.AdvanceToken() // consume "or"
			if orNode.FallbackStmt, err = p.ParseStmt(); err != nil {
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
	require.IsType(t, &jsextended.VarStmt{}, result.Stmts[0])
	stmt := result.Stmts[0].(*jsextended.VarStmt)
	require.IsType(t, &orExpr{}, stmt.Value)
	orVal := stmt.Value.(*orExpr)
	require.IsType(t, &js.ExprStmt{}, orVal.FallbackStmt)
	fallback := orVal.FallbackStmt.(*js.ExprStmt)
	require.IsType(t, &js.CallExpr{}, fallback.Expr)
	expr := fallback.Expr.(*js.CallExpr)
	require.IsType(t, &js.Variable{}, expr.Callee)
	funcName := expr.Callee.(*js.Variable)
	assert.Equal(t, "exit", funcName.Token.Literal)
}

type notBitwiseExpr struct {
	ast.BaseExpr
	Operator token.Token
	Value    ast.Expr
}

func TestUseUnaryParser(t *testing.T) {
	notBitwise := token.RegisterType("~")
	token.RegisterUnaryType(notBitwise)
	input := "1 + ~7"
	b := xjs.PluginBuilder()
	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
		if sc.CurrentChar() == '~' {
			sc.AdvanceChar()
			return token.Token{Type: notBitwise, Literal: "~"}, nil
		}
		return next()
	})
	b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (node ast.Expr, err error) {
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
	require.IsType(t, &js.Literal{}, binNode.Left)
	require.Equal(t, "1", binNode.Left.(*js.Literal).Value.Literal)
	require.Equal(t, token.PLUS, binNode.Op.Type)
	require.IsType(t, &notBitwiseExpr{}, binNode.Right)
	rightNode := binNode.Right.(*notBitwiseExpr)
	require.IsType(t, &js.Literal{}, rightNode.Value)
	require.Equal(t, "7", rightNode.Value.(*js.Literal).Value.Literal)
	require.Equal(t, notBitwise, rightNode.Operator.Type)
}

type powExpr struct {
	ast.BaseExpr
	LeftValue  ast.Expr
	Operator   token.Token
	RightValue ast.Expr
}

func TestUseBinaryParser(t *testing.T) {
	powType := token.RegisterType("^")
	token.RegisterBinaryType(powType, token.MULTIPLY.Precedence()+1)
	input := "1+5^2"
	b := xjs.PluginBuilder()
	b.UseScanner(func(s *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
		if s.CurrentChar() == '^' {
			s.AdvanceChar()
			return token.Token{Type: powType, Literal: "^"}, nil
		}
		return next()
	})
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(ast.Expr) (ast.Expr, error)) (node ast.Expr, err error) {
		if p.CurrentToken.Type == powType {
			powNode := &powExpr{LeftValue: left, Operator: p.CurrentToken}
			p.AdvanceToken() // consume ^
			if powNode.RightValue, err = js.ParseRightExpr(p, powType.Precedence()); err != nil {
				return
			}
			node = powNode
			return
		}
		return next(left)
	})
	p := b.Build([]byte(input))
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.IsType(t, &js.BinaryExpr{}, result)
	binNode := result.(*js.BinaryExpr)
	require.IsType(t, &js.Literal{}, binNode.Left)
	require.Equal(t, "1", binNode.Left.(*js.Literal).Value.Literal)
	require.Equal(t, token.PLUS, binNode.Op.Type)
	require.IsType(t, &powExpr{}, binNode.Right)
	powNode := binNode.Right.(*powExpr)
	require.IsType(t, &js.Literal{}, powNode.LeftValue)
	require.Equal(t, "5", powNode.LeftValue.(*js.Literal).Value.Literal)
	require.Equal(t, powType, powNode.Operator.Type)
	require.IsType(t, &js.Literal{}, powNode.RightValue)
	require.Equal(t, "2", powNode.RightValue.(*js.Literal).Value.Literal)
}

type factorialExpr struct {
	ast.BaseExpr
	Operator token.Token
	Value    ast.Expr
}

func TestUseBinaryParser_postfix(t *testing.T) {
	facTyp := token.RegisterType("!")
	token.RegisterBinaryType(facTyp, -1)
	input := "5! + 1"
	b := xjs.PluginBuilder()
	b.UseScanner(func(s *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
		if s.CurrentChar() == '!' {
			s.AdvanceChar()
			return token.Token{Type: facTyp, Literal: "!"}, nil
		}
		return next()
	})
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (node ast.Expr, err error) {
		if p.CurrentToken.Type == facTyp {
			left = &factorialExpr{Operator: p.CurrentToken, Value: left}
			p.AdvanceToken() // consume !
		}
		return next(left)
	})
	p := b.Build([]byte(input))
	result, err := p.ParseExpr()
	if err != nil {
		t.Fatal(err)
	}
	// check the result
	require.IsType(t, &js.BinaryExpr{}, result)
	binNode := result.(*js.BinaryExpr)
	require.IsType(t, &factorialExpr{}, binNode.Left)
	leftNode := binNode.Left.(*factorialExpr)
	require.IsType(t, &js.Literal{}, leftNode.Value)
	left := leftNode.Value.(*js.Literal)
	require.Equal(t, "5", left.Value.Literal)
	require.Equal(t, facTyp, leftNode.Operator.Type)
	require.IsType(t, &js.Literal{}, binNode.Right)
	rigth := binNode.Right.(*js.Literal)
	require.Equal(t, "1", rigth.Value.Literal)
}
