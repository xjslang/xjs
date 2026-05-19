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
	switch sc.CurrentChar {
	// operators
	case '=':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return token.Token{Type: token.EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.ASSIGN, Literal: string(c1)}
	case '!':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return token.Token{Type: token.NOT_EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.NOT, Literal: string(c1)}
	case '<':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return token.Token{Type: token.LTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.LT, Literal: string(c1)}
	case '>':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return token.Token{Type: token.GTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.GT, Literal: string(c1)}
	// maths operators
	case '+':
		sc.AdvanceChar()
		return token.Token{Type: token.PLUS, Literal: token.PLUS.String()}
	case '-':
		sc.AdvanceChar()
		return token.Token{Type: token.MINUS, Literal: token.MINUS.String()}
	case '*':
		sc.AdvanceChar()
		return token.Token{Type: token.MULTIPLY, Literal: token.MULTIPLY.String()}
	case '%':
		sc.AdvanceChar()
		return token.Token{Type: token.MODULO, Literal: token.MODULO.String()}
	// divide operator and comments
	case '/':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '/' {
			comment := sc.consumeLineComment()
			return token.Token{Type: token.LINE_COMMENT, Literal: comment}
		}
		if sc.CurrentChar == '*' {
			comment, typ := sc.consumeBlockComment()
			return token.Token{Type: typ, Literal: comment}
		}
		return token.Token{Type: token.DIVIDE, Literal: string(c1)}
	case '\'', '"':
		lit, typ := sc.consumeString(sc.CurrentChar)
		return token.Token{Type: typ, Literal: lit}
	// delimiters
	case ',':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return token.Token{Type: token.COMMA, Literal: string(c)}
	case ';':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return token.Token{Type: token.SEMICOLON, Literal: string(c)}
	case '(':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return token.Token{Type: token.LPAREN, Literal: string(c)}
	case ')':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return token.Token{Type: token.RPAREN, Literal: string(c)}
	case '{':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return token.Token{Type: token.LBRACE, Literal: string(c)}
	case '}':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return token.Token{Type: token.RBRACE, Literal: string(c)}
	case '\r':
		sc.AdvanceChar()
		if sc.CurrentChar == '\n' {
			sc.AdvanceChar()
			return token.Token{Type: token.NEWLINE, Literal: "\r\n"}
		}
		return token.Token{Type: token.NEWLINE, Literal: "\r"}
	case '\n':
		sc.AdvanceChar()
		return token.Token{Type: token.NEWLINE, Literal: "\n"}
	default:
		if isLetter(sc.CurrentChar) {
			lit := sc.consumeIdentifier()
			typ := lookup(lit)
			return token.Token{Type: typ, Literal: lit}
		} else if isDigit(sc.CurrentChar) {
			lit := sc.consumeNumber()
			return token.Token{Type: token.NUMBER, Literal: lit}
		} else if sc.CurrentChar == utf8.RuneError {
			c := sc.CurrentChar
			sc.AdvanceChar()
			return token.Token{Type: token.ILLEGAL, Literal: string(c)}
		} else if sc.CurrentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	c := sc.CurrentChar
	sc.AdvanceChar()
	return token.Token{Type: token.UNKNOWN, Literal: string(c)}
}

func lookup(lit string) token.Type {
	switch lit {
	case "true", "false":
		return token.BOOLEAN
	case "let":
		return token.LET
	case "function":
		return token.FUNCTION
	}
	return token.IDENT
}
