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
		c := l.CurrentChar
		l.AdvanceChar()
		return token.Token{Type: token.SEMICOLON, Literal: string(c)}
	case '=':
		c1 := l.CurrentChar
		l.AdvanceChar()
		if l.CurrentChar == '=' {
			c2 := l.CurrentChar
			l.AdvanceChar()
			return token.Token{Type: token.EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.ASSIGN, Literal: string(c1)}
	case '!':
		c1 := l.CurrentChar
		l.AdvanceChar()
		if l.CurrentChar == '=' {
			c2 := l.CurrentChar
			l.AdvanceChar()
			return token.Token{Type: token.NOT_EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.NOT, Literal: string(c1)}
	case '<':
		c1 := l.CurrentChar
		l.AdvanceChar()
		if l.CurrentChar == '=' {
			c2 := l.CurrentChar
			l.AdvanceChar()
			return token.Token{Type: token.LTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.LT, Literal: string(c1)}
	case '>':
		c1 := l.CurrentChar
		l.AdvanceChar()
		if l.CurrentChar == '=' {
			c2 := l.CurrentChar
			l.AdvanceChar()
			return token.Token{Type: token.GTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.GT, Literal: string(c1)}
	case '/':
		c1 := l.CurrentChar
		l.AdvanceChar()
		if l.CurrentChar == '/' {
			l.AdvanceChar()
			comment := l.parseSinglelineComment()
			return token.Token{Type: token.LCOMMENT, Literal: comment}
		}
		if l.CurrentChar == '*' {
			l.AdvanceChar()
			comment, typ := l.parseMultilineComment()
			return token.Token{Type: typ, Literal: comment}
		}

		return token.Token{Type: token.DIVIDE, Literal: string(c1)}
	case '\'', '"':
		lit, typ := l.parseString(l.CurrentChar)
		return token.Token{Type: typ, Literal: lit}
	case '(':
		c := l.CurrentChar
		l.AdvanceChar()
		return token.Token{Type: token.LPAREN, Literal: string(c)}
	case ')':
		c := l.CurrentChar
		l.AdvanceChar()
		return token.Token{Type: token.RPAREN, Literal: string(c)}
	case '{':
		c := l.CurrentChar
		l.AdvanceChar()
		return token.Token{Type: token.LBRACE, Literal: string(c)}
	case '}':
		c := l.CurrentChar
		l.AdvanceChar()
		return token.Token{Type: token.RBRACE, Literal: string(c)}
	case '\r':
		l.AdvanceChar()
		if l.CurrentChar == '\n' {
			l.AdvanceChar()
		}
		return token.Token{Type: token.NEWLINE, Literal: ""}
	case '\n':
		l.AdvanceChar()
		return token.Token{Type: token.NEWLINE, Literal: ""}
	default:
		if isLetter(l.CurrentChar) {
			lit := l.parseIdentifier()
			typ := token.Lookup(lit)
			return token.Token{Type: typ, Literal: lit}
		} else if isDigit(l.CurrentChar) {
			lit := l.parseNumber()
			return token.Token{Type: token.NUMBER, Literal: lit}
		} else if l.CurrentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	c := l.CurrentChar
	l.AdvanceChar()
	return token.Token{Type: token.UNKNOWN, Literal: string(c)}
}
