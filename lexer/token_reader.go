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
	case ';':
		lit := l.consumeChar()
		return token.Token{Type: token.SEMI, Literal: lit}
	case '=':
		if l.PeekChar == '=' {
			lit := l.consumeChars(2)
			return token.Token{Type: token.EQ, Literal: lit}
		} else {
			lit := l.consumeChar()
			return token.Token{Type: token.ASSIGN, Literal: lit}
		}
	case '!':
		if l.PeekChar == '=' {
			lit := l.consumeChars(2)
			return token.Token{Type: token.NOT_EQ, Literal: lit}
		} else {
			lit := l.consumeChar()
			return token.Token{Type: token.NOT, Literal: lit}
		}
	case '<':
		if l.PeekChar == '=' {
			lit := l.consumeChars(2)
			return token.Token{Type: token.LOWER_OR_EQ, Literal: lit}
		} else {
			lit := l.consumeChar()
			return token.Token{Type: token.LOWER, Literal: lit}
		}
	case '>':
		if l.PeekChar == '=' {
			lit := l.consumeChars(2)
			return token.Token{Type: token.GREATER_OR_EQ, Literal: lit}
		} else {
			lit := l.consumeChar()
			return token.Token{Type: token.GREATER, Literal: lit}
		}
	case '\'', '"':
		lit, typ := l.consumeString(l.CurrentChar)
		return token.Token{Type: typ, Literal: lit}
	case '(':
		lit := l.consumeChar()
		return token.Token{Type: token.LPAREN, Literal: lit}
	case ')':
		lit := l.consumeChar()
		return token.Token{Type: token.RPAREN, Literal: lit}
	case '{':
		lit := l.consumeChar()
		return token.Token{Type: token.LBRACE, Literal: lit}
	case '}':
		lit := l.consumeChar()
		return token.Token{Type: token.RBRACE, Literal: lit}
	default:
		if isLetter(l.CurrentChar) {
			lit := l.consumeIdentifier()
			typ := token.Lookup(lit)
			return token.Token{Type: typ, Literal: lit}
		} else if isDigit(l.CurrentChar) {
			lit := l.consumeNumber()
			return token.Token{Type: token.NUMBER, Literal: lit}
		} else if l.CurrentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	lit := l.consumeChar()
	return token.Token{Type: token.UNKNOWN, Literal: lit}
}
