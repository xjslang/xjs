package lexer

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

func (l *Lexer) UseTokenizer(tokenizer func(l *Lexer, next func() token.Token) token.Token) {
	next := l.tokenizer
	if next == nil {
		next = defaultTokenizer
	}
	l.tokenizer = func(l *Lexer) token.Token {
		return tokenizer(l, func() token.Token {
			return next(l)
		})
	}
}

func defaultTokenizer(l *Lexer) token.Token {
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
			comment := l.consumeLineComment()
			return token.Token{Type: token.LINE_COMMENT, Literal: comment}
		}
		if l.CurrentChar == '*' {
			comment, typ := l.consumeBlockComment()
			return token.Token{Type: typ, Literal: comment}
		}

		return token.Token{Type: token.DIVIDE, Literal: string(c1)}
	case '\'', '"':
		lit, typ := l.consumeString(l.CurrentChar)
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
			lit := l.consumeIdentifier()
			typ := token.Lookup(lit)
			return token.Token{Type: typ, Literal: lit}
		} else if isDigit(l.CurrentChar) {
			lit := l.consumeNumber()
			return token.Token{Type: token.NUMBER, Literal: lit}
		} else if l.CurrentChar == utf8.RuneError {
			c := l.CurrentChar
			l.AdvanceChar()
			return token.Token{Type: token.ILLEGAL, Literal: string(c)}
		} else if l.CurrentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	c := l.CurrentChar
	l.AdvanceChar()
	return token.Token{Type: token.UNKNOWN, Literal: string(c)}
}
