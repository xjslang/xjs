package scanner

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

func (s *Scanner) UseScanner(scanner func(s *Scanner, next func() token.Token) token.Token) {
	next := s.scanner
	if next == nil {
		next = defaultScanner
	}
	s.scanner = func(s *Scanner) token.Token {
		return scanner(s, func() token.Token {
			return next(s)
		})
	}
}

func defaultScanner(s *Scanner) token.Token {
	switch s.currentChar {
	// operators
	case '=':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.ASSIGN, Literal: string(c1)}
	case '!':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.NOT_EQ, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.NOT, Literal: string(c1)}
	case '<':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.LTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.LT, Literal: string(c1)}
	case '>':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.GTE, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.GT, Literal: string(c1)}
	case '|':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '|' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.OR, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.UNKNOWN, Literal: string(c1)}
	case '&':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '&' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.AND, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.UNKNOWN, Literal: string(c1)}
	// maths operators
	case '+':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '+' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.INCREMENT, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.PLUS, Literal: string(c1)}
	case '-':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '-' {
			c2 := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.DECREMENT, Literal: string([]rune{c1, c2})}
		}
		return token.Token{Type: token.MINUS, Literal: string(c1)}
	case '*':
		s.AdvanceChar()
		return token.Token{Type: token.MULTIPLY, Literal: token.MULTIPLY.String()}
	case '%':
		s.AdvanceChar()
		return token.Token{Type: token.MODULO, Literal: token.MODULO.String()}
	// divide operator and comments
	case '/':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '/' {
			comment := scanLineComment(s)
			return token.Token{Type: token.LINE_COMMENT, Literal: comment}
		}
		if s.currentChar == '*' {
			comment, typ := scanBlockComment(s)
			return token.Token{Type: typ, Literal: comment}
		}
		return token.Token{Type: token.DIVIDE, Literal: string(c1)}
	case '\'', '"':
		lit, typ := scanString(s, s.currentChar)
		return token.Token{Type: typ, Literal: lit}
	// delimiters
	case ',':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.COMMA, Literal: string(c)}
	case '.':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.DOT, Literal: string(c)}
	case ';':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.SEMICOLON, Literal: string(c)}
	case '(':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.LPAREN, Literal: string(c)}
	case ')':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.RPAREN, Literal: string(c)}
	case '{':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.LBRACE, Literal: string(c)}
	case '}':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.RBRACE, Literal: string(c)}
	case '[':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.LBRACKET, Literal: string(c)}
	case ']':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.RBRACKET, Literal: string(c)}
	case ':':
		c := s.currentChar
		s.AdvanceChar()
		return token.Token{Type: token.COLON, Literal: string(c)}
	case '\r':
		s.AdvanceChar()
		if s.currentChar == '\n' {
			s.AdvanceChar()
			return token.Token{Type: token.NEWLINE, Literal: "\r\n"}
		}
		return token.Token{Type: token.NEWLINE, Literal: "\r"}
	case '\n':
		s.AdvanceChar()
		return token.Token{Type: token.NEWLINE, Literal: "\n"}
	default:
		if isLetter(s.currentChar) {
			lit := scanIdentifier(s)
			return token.Token{Type: token.IDENT, Literal: lit}
		} else if isDigit(s.currentChar) {
			lit, typ := scanNumber(s)
			return token.Token{Type: typ, Literal: lit}
		} else if s.currentChar == utf8.RuneError {
			c := s.currentChar
			s.AdvanceChar()
			return token.Token{Type: token.ILLEGAL, Literal: string(c)}
		} else if s.currentChar == eof {
			return token.Token{Type: token.EOF, Literal: ""}
		}
	}

	c := s.currentChar
	s.AdvanceChar()
	return token.Token{Type: token.UNKNOWN, Literal: string(c)}
}
