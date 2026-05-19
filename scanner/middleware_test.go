package scanner_test

import (
	"testing"

	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func TestUseScanner(t *testing.T) {
	sc := &scanner.Scanner{}
	powType := token.RegisterType("**")
	sc.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		if sc.CurrentChar == '*' && sc.PeekChar() == '*' {
			// consume **
			sc.AdvanceChar()
			sc.AdvanceChar()
			return token.Token{Type: powType, Literal: powType.String()}
		}
		return next()
	})
	sc.Init([]byte("125 ** 2"))
	assertLexerTokens(t, sc, []token.Token{
		{Type: token.NUMBER, Literal: "125"},
		{Type: powType, Literal: "**"},
		{Type: token.NUMBER, Literal: "2"},
		{Type: token.EOF},
	})
}
