package scanner

import (
	"unicode/utf8"
)

func (sc *Scanner) UseScanner(scanner func(sc *Scanner, next func() Token) Token) {
	next := sc.scanner
	if next == nil {
		next = defaultScanner
	}
	sc.scanner = func(sc *Scanner) Token {
		return scanner(sc, func() Token {
			return next(sc)
		})
	}
}

func defaultScanner(sc *Scanner) Token {
	switch sc.CurrentChar {
	// operators
	case '=':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return Token{Type: EQ, Literal: string([]rune{c1, c2})}
		}
		return Token{Type: ASSIGN, Literal: string(c1)}
	case '!':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return Token{Type: NOT_EQ, Literal: string([]rune{c1, c2})}
		}
		return Token{Type: NOT, Literal: string(c1)}
	case '<':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return Token{Type: LTE, Literal: string([]rune{c1, c2})}
		}
		return Token{Type: LT, Literal: string(c1)}
	case '>':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '=' {
			c2 := sc.CurrentChar
			sc.AdvanceChar()
			return Token{Type: GTE, Literal: string([]rune{c1, c2})}
		}
		return Token{Type: GT, Literal: string(c1)}
	// maths operators
	case '+':
		sc.AdvanceChar()
		return Token{Type: PLUS, Literal: PLUS.String()}
	case '-':
		sc.AdvanceChar()
		return Token{Type: MINUS, Literal: MINUS.String()}
	case '*':
		sc.AdvanceChar()
		return Token{Type: MULTIPLY, Literal: MULTIPLY.String()}
	case '%':
		sc.AdvanceChar()
		return Token{Type: MODULO, Literal: MODULO.String()}
	// divide operator and comments
	case '/':
		c1 := sc.CurrentChar
		sc.AdvanceChar()
		if sc.CurrentChar == '/' {
			comment := sc.consumeLineComment()
			return Token{Type: LINE_COMMENT, Literal: comment}
		}
		if sc.CurrentChar == '*' {
			comment, typ := sc.consumeBlockComment()
			return Token{Type: typ, Literal: comment}
		}
		return Token{Type: DIVIDE, Literal: string(c1)}
	case '\'', '"':
		lit, typ := sc.consumeString(sc.CurrentChar)
		return Token{Type: typ, Literal: lit}
	// delimiters
	case ';':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return Token{Type: SEMICOLON, Literal: string(c)}
	case '(':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return Token{Type: LPAREN, Literal: string(c)}
	case ')':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return Token{Type: RPAREN, Literal: string(c)}
	case '{':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return Token{Type: LBRACE, Literal: string(c)}
	case '}':
		c := sc.CurrentChar
		sc.AdvanceChar()
		return Token{Type: RBRACE, Literal: string(c)}
	case '\r':
		sc.AdvanceChar()
		if sc.CurrentChar == '\n' {
			sc.AdvanceChar()
		}
		return Token{Type: NEWLINE, Literal: ""}
	case '\n':
		sc.AdvanceChar()
		return Token{Type: NEWLINE, Literal: ""}
	default:
		if isLetter(sc.CurrentChar) {
			lit := sc.consumeIdentifier()
			typ := lookup(lit)
			return Token{Type: typ, Literal: lit}
		} else if isDigit(sc.CurrentChar) {
			lit := sc.consumeNumber()
			return Token{Type: NUMBER, Literal: lit}
		} else if sc.CurrentChar == utf8.RuneError {
			c := sc.CurrentChar
			sc.AdvanceChar()
			return Token{Type: ILLEGAL, Literal: string(c)}
		} else if sc.CurrentChar == eof {
			return Token{Type: EOF, Literal: ""}
		}
	}

	c := sc.CurrentChar
	sc.AdvanceChar()
	return Token{Type: UNKNOWN, Literal: string(c)}
}

func lookup(lit string) Kind {
	switch lit {
	case "true", "false":
		return BOOLEAN
	case "let":
		return LET
	case "function":
		return FUNCTION
	}
	return IDENT
}
