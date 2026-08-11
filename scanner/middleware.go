package scanner

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

// TODO: In the parameter type for `scanner`, the argument is named `s *Scanner`, which is easy to confuse with the method receiver `s *Scanner`. Renaming it to `sc` (or omitting parameter names in the func type) improves readability without changing the type.
func (s *Scanner) useScanner(scanner func(s *Scanner, next func(*Scanner) (token.Token, error)) (token.Token, error)) {
	next := s.scanner
	if next == nil {
		next = defaultScanner
	}
	s.scanner = func(s *Scanner) (token.Token, error) {
		return scanner(s, next)
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
		tok = token.Token{Type: token.MULTIPLY, Literal: token.MULTIPLY.Literal()}
	case '%':
		s.AdvanceChar()
		tok = token.Token{Type: token.MODULO, Literal: token.MODULO.Literal()}
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
				tok.Type = token.ILLEGAL
				return
			}
		default:
			tok = token.Token{Type: token.DIVIDE, Literal: string(c)}
		}
	// delimiters
	case '\'', '"':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.STRING}
		if tok.Literal, err = ScanString(s, c); err != nil {
			tok.Type = token.ILLEGAL
			return
		}
	case '`':
		s.AdvanceChar()
		tok = token.Token{Type: token.STRING}
		if tok.Literal, err = ScanRawString(s); err != nil {
			tok.Type = token.ILLEGAL
			return
		}
	case ',':
		c := s.currentChar
		s.AdvanceChar()
		tok = token.Token{Type: token.COMMA, Literal: string(c)}
	case '.':
		c := s.currentChar
		if IsDigit(s.PeekChar()) {
			tok.Type = token.NUMBER
			if tok.Literal, err = ScanNumber(s); err != nil {
				tok.Type = token.ILLEGAL
				return
			}
		} else {
			s.AdvanceChar()
			tok = token.Token{Type: token.DOT, Literal: string(c)}
		}
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
			tok.Type = token.NUMBER
			if tok.Literal, err = ScanNumber(s); err != nil {
				tok.Type = token.ILLEGAL
				return
			}
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
