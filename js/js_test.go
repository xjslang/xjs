package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func createParser(t *testing.T, input string) *parser.Parser {
	t.Helper()
	s := &scanner.Scanner{}
	s.Init([]byte(input))
	p := &parser.Parser{}
	p.Init(s)
	return p
}

func TestExpectSemi(t *testing.T) {
	t.Run("no scope", func(t *testing.T) {
		input := "hello; there"
		p := createParser(t, input)
		// expect ident
		tok, err := p.Expect(token.IDENT)
		require.NoError(t, err)
		require.Equal(t, "hello", tok.Literal)
		// expect semi
		tok, err = ExpectSemi(p)
		require.NoError(t, err)
		// semicolon is consumed; expect identifier
		tok, err = p.Expect(token.IDENT)
		assert.NoError(t, err)
		assert.Equal(t, "there", tok.Literal)
		// expect EOF
		tok, err = p.Expect(token.EOF)
		require.NoError(t, err)
	})

	t.Run("block scope", func(t *testing.T) {
		input := "hello} there"
		p := createParser(t, input)
		p.EnterScope(blockScope)
		defer p.ExitScope(blockScope)
		// expect ident
		tok, err := p.Expect(token.IDENT)
		require.NoError(t, err)
		require.Equal(t, "hello", tok.Literal)
		// '}' is a semicolon delimiter in a "block scope"
		tok, err = ExpectSemi(p)
		require.NoError(t, err)
		// '}' is NOT consumed in a "block scope"
		tok, err = p.Expect(token.RBRACE)
		require.NoError(t, err)
		require.Equal(t, "}", tok.Literal)
		// expect ident
		tok, err = p.Expect(token.IDENT)
		require.NoError(t, err)
		require.Equal(t, "there", tok.Literal)
		// expect EOF
		tok, err = p.Expect(token.EOF)
		require.NoError(t, err)
	})

	t.Run("for scope", func(t *testing.T) {
		input := "hello; there)"
		p := createParser(t, input)
		p.EnterScope(forHeaderScope)
		defer p.ExitScope(forHeaderScope)
		// expect ident
		tok, err := p.Expect(token.IDENT)
		require.NoError(t, err)
		require.Equal(t, "hello", tok.Literal)
		// expect semi
		tok, err = ExpectSemi(p)
		require.NoError(t, err)
		// ';' is not consumed in a "for scope"
		tok, err = p.Expect(token.SEMICOLON)
		require.NoError(t, err)
		require.Equal(t, ";", tok.Literal)
		// expect ident
		tok, err = p.Expect(token.IDENT)
		require.NoError(t, err)
		require.Equal(t, "there", tok.Literal)
		// ')' is a semicolon delimiter in a "for scope"
		tok, err = ExpectSemi(p)
		require.NoError(t, err)
		// ')' is not consumed in a "for scope"
		tok, err = p.Expect(token.RPAREN)
		require.NoError(t, err)
		require.Equal(t, ")", tok.Literal)
		// expect EOF
		tok, err = p.Expect(token.EOF)
		require.NoError(t, err)
	})
}
