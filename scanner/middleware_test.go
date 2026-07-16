package scanner_test

import (
	"testing"

	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func TestUseScanner(t *testing.T) {
	powType := token.RegisterType("**")
	sc := scanner.NewBuilder().
		UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
			if sc.CurrentChar() == '*' && sc.PeekChar() == '*' {
				// consume **
				sc.AdvanceChar()
				sc.AdvanceChar()
				return token.Token{Type: powType, Literal: powType.String()}, nil
			}
			return next()
		}).
		Build([]byte("5 ** 2"))
	assertLexerTokens(t, sc, []token.Token{
		{Type: token.NUMBER, Literal: "5"},
		{Type: powType, Literal: "**"},
		{Type: token.NUMBER, Literal: "2"},
		{Type: token.EOF},
	})
}
