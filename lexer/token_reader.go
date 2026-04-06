package lexer

import (
	"github.com/xjslang/xjs/token"
)

func (l *Lexer) UseTokenReader(reader func(l *Lexer, next func() token.Token) token.Token) {
	next := l.tokenReader
	l.tokenReader = func(l *Lexer) token.Token {
		return reader(l, func() token.Token {
			return next(l)
		})
	}
}

func defaultTokenReader(l *Lexer) token.Token {
	switch l.CurrentChar {
	case '=':
		if l.PeekChar == '=' {
			lit := l.readChars(2)
			return token.Token{Type: token.EQ, Literal: lit}
		} else {
			lit := l.readChar()
			return token.Token{Type: token.ASSIGN, Literal: lit}
		}
	case '!':
		if l.PeekChar == '=' {
			lit := l.readChars(2)
			return token.Token{Type: token.NOT_EQ, Literal: lit}
		} else {
			lit := l.readChar()
			return token.Token{Type: token.NOT, Literal: lit}
		}
	case '<':
		if l.PeekChar == '=' {
			lit := l.readChars(2)
			return token.Token{Type: token.LOWER_OR_EQ, Literal: lit}
		} else {
			lit := l.readChar()
			return token.Token{Type: token.LOWER, Literal: lit}
		}
	case '>':
		if l.PeekChar == '=' {
			lit := l.readChars(2)
			return token.Token{Type: token.GREATER_OR_EQ, Literal: lit}
		} else {
			lit := l.readChar()
			return token.Token{Type: token.GREATER, Literal: lit}
		}
	default:
		if isLetter(l.CurrentChar) {
			lit := l.readIden()
			return token.Token{Type: token.IDENT, Literal: lit}
		} else if isDigit(l.CurrentChar) {
			lit := l.readNumber()
			return token.Token{Type: token.NUMBER, Literal: lit}
		} else if l.CurrentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	lit := l.readChar()
	return token.Token{Type: token.UNKNOWN, Literal: lit}
}
