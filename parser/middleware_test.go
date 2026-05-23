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
