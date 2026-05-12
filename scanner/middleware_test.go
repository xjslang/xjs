package scanner_test

import (
	"testing"

	"github.com/xjslang/xjs/scanner"
)

func TestUseScanner(t *testing.T) {
	sc := &scanner.Scanner{}
	powType := scanner.RegisterKind("**")
	sc.UseScanner(func(sc *scanner.Scanner, next func() scanner.Token) scanner.Token {
		if sc.CurrentChar == '*' && sc.PeekChar() == '*' {
			// consume **
			sc.AdvanceChar()
			sc.AdvanceChar()
			return scanner.Token{Type: powType, Literal: powType.String()}
		}
		return next()
	})
	sc.Init([]byte("125 ** 2"))
	assertLexerTokens(t, sc, []scanner.Token{
		{Type: scanner.NUMBER, Literal: "125"},
		{Type: powType, Literal: "**"},
		{Type: scanner.NUMBER, Literal: "2"},
		{Type: scanner.EOF},
	})
}
