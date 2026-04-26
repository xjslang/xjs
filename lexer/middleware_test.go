package lexer

import (
	"testing"

	"github.com/xjslang/xjs/token"
)

func TestUseTokenizer(t *testing.T) {
	l := &Lexer{}
	powType := token.RegisterType("**")
	l.UseTokenizer(func(l *Lexer, next func() token.Token) token.Token {
		if l.CurrentChar == '*' && l.PeekChar() == '*' {
			// consume **
			l.AdvanceChar()
			l.AdvanceChar()
			return token.Token{Type: powType, Literal: powType.String()}
		}
		return next()
	})
	l.Init([]byte("125 ** 2"))
	assertLexerTokens(t, l, []token.Token{
		{Type: token.NUMBER, Literal: "125"},
		{Type: powType, Literal: "**"},
		{Type: token.NUMBER, Literal: "2"},
		{Type: token.EOF},
	})
}
