package scanner_test

import (
	"testing"

	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func TestUseScanner(t *testing.T) {
	sc := &scanner.Scanner{}
	powType := token.RegisterType("**")
	sc.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (token.Token, error) {
		if sc.CurrentChar() == '*' && sc.PeekChar() == '*' {
			// consume **
			sc.AdvanceChar()
			sc.AdvanceChar()
			return token.Token{Type: powType, Literal: powType.String()}, nil
		}
		return next()
	})
	sc.Init([]byte("5 ** 2"))
	assertLexerTokens(t, sc, []token.Token{
		{Type: token.DIGIT, Literal: "5"},
		{Type: powType, Literal: "**"},
		{Type: token.DIGIT, Literal: "2"},
		{Type: token.EOF},
	})
}
