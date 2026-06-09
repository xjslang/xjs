package scanner

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

func (sc *Scanner) UseScanner(scanner func(sc *Scanner, next func() token.Token) token.Token) {
	next := sc.scanner
	if next == nil {
		next = defaultScanner
	}
	sc.scanner = func(sc *Scanner) token.Token {
		return scanner(sc, func() token.Token {
			return next(sc)
		})
	}
}

func defaultScanner(sc *Scanner) token.Token {
	switch sc.currentChar {
	// operators
	case '=':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '=' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.ASSIGN, Literal: string(c1)}
	case '!':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '=' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.NOT_EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.NOT, Literal: string(c1)}
	case '<':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '=' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.LTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.LT, Literal: string(c1)}
	case '>':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '=' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.GTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.GT, Literal: string(c1)}
	case '|':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '|' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.OR, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.UNKNOWN, Literal: string(c1)}
	case '&':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '&' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.AND, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.UNKNOWN, Literal: string(c1)}
	// maths operators
	case '+':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '+' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.INCREMENT, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.PLUS, Literal: string(c1)}
	case '-':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '-' {
			c2 := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.DECREMENT, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.MINUS, Literal: string(c1)}
	case '*':
		sc.AdvanceChar()
		return token.Token{Type: token.MULTIPLY, Literal: token.MULTIPLY.String()}
	case '%':
		sc.AdvanceChar()
		return token.Token{Type: token.MODULO, Literal: token.MODULO.String()}
	// divide operator and comments
	case '/':
		c1 := sc.currentChar
		sc.AdvanceChar()
		if sc.currentChar == '/' {
			comment := scanLineComment(sc)
			return token.Token{Type: token.LINE_COMMENT, Literal: comment}
		}
		if sc.currentChar == '*' {
			comment, typ := scanBlockComment(sc)
			return token.Token{Type: typ, Literal: comment}
		}
		return token.Token{Type: token.DIVIDE, Literal: string(c1)}
	case '\'', '"':
		lit, typ := scanString(sc, sc.currentChar)
		return token.Token{Type: typ, Literal: lit}
	// delimiters
	case ',':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.COMMA, Literal: string(c)}
	case '.':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.DOT, Literal: string(c)}
	case ';':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.SEMICOLON, Literal: string(c)}
	case '(':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.LPAREN, Literal: string(c)}
	case ')':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.RPAREN, Literal: string(c)}
	case '{':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.LBRACE, Literal: string(c)}
	case '}':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.RBRACE, Literal: string(c)}
	case '[':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.LBRACKET, Literal: string(c)}
	case ']':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.RBRACKET, Literal: string(c)}
	case ':':
		c := sc.currentChar
		sc.AdvanceChar()
		return token.Token{Type: token.COLON, Literal: string(c)}
	case '\r':
		sc.AdvanceChar()
		if sc.currentChar == '\n' {
			sc.AdvanceChar()
			return token.Token{Type: token.NEWLINE, Literal: "\r\n"}
		}
		return token.Token{Type: token.NEWLINE, Literal: "\r"}
	case '\n':
		sc.AdvanceChar()
		return token.Token{Type: token.NEWLINE, Literal: "\n"}
	default:
		if isLetter(sc.currentChar) {
			lit := scanIdentifier(sc)
			return token.Token{Type: token.IDENT, Literal: lit}
		} else if isDigit(sc.currentChar) {
			lit, typ := scanNumber(sc)
			return token.Token{Type: typ, Literal: lit}
		} else if sc.currentChar == utf8.RuneError {
			c := sc.currentChar
			sc.AdvanceChar()
			return token.Token{Type: token.ILLEGAL, Literal: string(c)}
		} else if sc.currentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	c := sc.currentChar
	sc.AdvanceChar()
	return token.Token{Type: token.UNKNOWN, Literal: string(c)}
}
