package scanner

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

func (s *Scanner) useScanner(scanner func(s *Scanner, next func() (token.Token, error)) (token.Token, error)) {
	next := s.scanner
	if next == nil {
		next = defaultScanner
	}
	s.scanner = func(s *Scanner) (token.Token, error) {
		return scanner(s, func() (token.Token, error) {
			return next(s)
		})
	}
}

func defaultScanner(s *Scanner) (tok token.Token, err error) {
	switch s.currentChar {
	// operators
	case '=':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.EQ, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(c1)}
		}
	case '!':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.NOT, Literal: string(c1)}
		}
	case '<':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.LTE, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(c1)}
		}
	case '>':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '=' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.GTE, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(c1)}
		}
	case '|':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '|' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.OR, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.UNKNOWN, Literal: string(c1)}
		}
	case '&':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '&' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.AND, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.UNKNOWN, Literal: string(c1)}
		}
	// maths operators
	case '+':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '+' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.INCREMENT, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(c1)}
		}
	case '-':
		c1 := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '-' {
			c2 := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.DECREMENT, Literal: string([]rune{c1, c2})}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(c1)}
		}
	case '*':
		s.AdvanceChar()
		tok = token.Token{Type: token.MULTIPLY, Literal: token.MULTIPLY.String()}
	case '%':
		s.AdvanceChar()
		tok = token.Token{Type: token.MODULO, Literal: token.MODULO.String()}
	// divide operator and comments
	case '/':
		c := s.currentChar
		s.AdvanceChar()
		switch s.currentChar {
		case '/':
			lit := ScanLineComment(s)
			tok = token.Token{Type: token.LINE_COMMENT, Literal: lit}
		case '*':
			tok = token.Token{Type: token.BLOCK_COMMENT}
			if tok.Literal, err = ScanBlockComment(s); err != nil {
				return
			}
		default:
			tok = token.Token{Type: token.DIVIDE, Literal: string(c)}
		}
	// delimiters
	case '\'', '"':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.QUOTE, Literal: string(c)}
	case ',':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.COMMA, Literal: string(c)}
	case '.':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.DOT, Literal: string(c)}
	case ';':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.SEMICOLON, Literal: string(c)}
	case '(':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.LPAREN, Literal: string(c)}
	case ')':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.RPAREN, Literal: string(c)}
	case '{':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.LBRACE, Literal: string(c)}
	case '}':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.RBRACE, Literal: string(c)}
	case '[':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.LBRACKET, Literal: string(c)}
		return
	case ']':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.RBRACKET, Literal: string(c)}
	case ':':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.COLON, Literal: string(c)}
	case '\r':
		s.AdvanceChar()
		if s.currentChar == '\n' {
			s.AdvanceChar()
			tok = token.Token{Type: token.NEWLINE, Literal: "\r\n"}
		} else {
			tok = token.Token{Type: token.NEWLINE, Literal: "\r"}
		}
	case '\n':
		s.AdvanceChar()
		tok = token.Token{Type: token.NEWLINE, Literal: "\n"}
	default:
		if IsLetter(s.currentChar) {
			lit := ScanIdentifier(s)
			tok = token.Token{Type: token.IDENT, Literal: lit}
		} else if IsDigit(s.currentChar) {
			c := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.DIGIT, Literal: string(c)}
		} else if s.currentChar == utf8.RuneError {
			c := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.ILLEGAL, Literal: string(c)}
		} else if s.currentChar == EOF {
			tok = token.Token{Type: token.EOF, Literal: ""}
		} else {
			c := s.currentChar
			s.AdvanceChar()
			tok = token.Token{Type: token.UNKNOWN, Literal: string(c)}
		}
	}
	return
}
