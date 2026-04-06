package lexer

import (
	"strings"

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
	sb := &strings.Builder{}
	switch l.CurrentChar {
	case '=':
		if l.PeekChar == '=' {
			l.consumeChars(sb, 2)
			return token.Token{Type: token.EQ, Literal: sb.String()}
		} else {
			l.consumeChar(sb)
			return token.Token{Type: token.ASSIGN, Literal: sb.String()}
		}
	case '!':
		if l.PeekChar == '=' {
			l.consumeChars(sb, 2)
			return token.Token{Type: token.NOT_EQ, Literal: sb.String()}
		} else {
			l.consumeChar(sb)
			return token.Token{Type: token.NOT, Literal: sb.String()}
		}
	case '<':
		if l.PeekChar == '=' {
			l.consumeChars(sb, 2)
			return token.Token{Type: token.LOWER_OR_EQ, Literal: sb.String()}
		} else {
			l.consumeChar(sb)
			return token.Token{Type: token.LOWER, Literal: sb.String()}
		}
	case '>':
		if l.PeekChar == '=' {
			l.consumeChars(sb, 2)
			return token.Token{Type: token.GREATER_OR_EQ, Literal: sb.String()}
		} else {
			l.consumeChar(sb)
			return token.Token{Type: token.GREATER, Literal: sb.String()}
		}
	default:
		if isLetter(l.CurrentChar) {
			l.consumeIdentifier(sb)
			return token.Token{Type: token.IDENT, Literal: sb.String()}
		} else if isDigit(l.CurrentChar) {
			l.consumeNumber(sb)
			return token.Token{Type: token.NUMBER, Literal: sb.String()}
		} else if l.CurrentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	l.consumeChar(sb)
	return token.Token{Type: token.UNKNOWN, Literal: sb.String()}
}
