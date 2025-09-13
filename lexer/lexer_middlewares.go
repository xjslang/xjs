package lexer

import "github.com/xjslang/xjs/token"

func (l *Lexer) UseTokenReader(reader func(l *Lexer, next func() token.Token) token.Token) {
	next := l.nextToken
	l.nextToken = func(l *Lexer) token.Token {
		l.skipWhitespace()
		return reader(l, func() token.Token {
			return next(l)
		})
	}
}
