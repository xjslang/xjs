package scanner

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

// TODO: In the parameter type for `scanner`, the argument is named `s *Scanner`, which is easy to confuse with the method receiver `s *Scanner`. Renaming it to `sc` (or omitting parameter names in the func type) improves readability without changing the type.
func (s *Scanner) useScanner(scanner func(s *Scanner, next func(*Scanner) token.Token) token.Token) {
	next := s.scanner
	if next == nil {
		next = defaultScanner
	}
	s.scanner = func(s *Scanner) token.Token {
		return scanner(s, next)
	}
}

func defaultScanner(s *Scanner) (tok token.Token) {
	ofs0 := s.offset - s.currentSize
	switch s.currentChar {
	// operators
	case '=':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type = token.EQ
		} else {
			tok.Type = token.ASSIGN
		}
	case '!':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type = token.NOT_EQ
		} else {
			tok.Type = token.NOT
		}
	case '<':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type = token.LTE
		} else {
			tok.Type = token.LT
		}
	case '>':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type = token.GTE
		} else {
			tok.Type = token.GT
		}
	case '|':
		s.AdvanceChar()
		if s.currentChar == '|' {
			s.AdvanceChar()
			tok.Type = token.OR
		} else {
			tok.Type = token.UNKNOWN
		}
	case '&':
		s.AdvanceChar()
		if s.currentChar == '&' {
			s.AdvanceChar()
			tok.Type = token.AND
		} else {
			tok.Type = token.UNKNOWN
		}
	// maths operators
	case '+':
		s.AdvanceChar()
		if s.currentChar == '+' {
			s.AdvanceChar()
			tok.Type = token.INCREMENT
		} else {
			tok.Type = token.PLUS
		}
	case '-':
		s.AdvanceChar()
		if s.currentChar == '-' {
			s.AdvanceChar()
			tok.Type = token.DECREMENT
		} else {
			tok.Type = token.MINUS
		}
	case '*':
		s.AdvanceChar()
		tok.Type = token.MULTIPLY
	case '%':
		s.AdvanceChar()
		tok.Type = token.MODULO
	// divide operator and comments
	case '/':
		s.AdvanceChar()
		tok.Type = token.DIVIDE
	// delimiters
	case '\'', '"':
		c := s.currentChar
		s.AdvanceChar()
		tok.Type = token.STRING
		if err := scanString(s, c); err != nil {
			tok.Type = token.ILLEGAL
		}
	case '`':
		s.AdvanceChar()
		tok.Type = token.STRING
		if err := scanRawString(s); err != nil {
			tok.Type = token.ILLEGAL
		}
	case ',':
		s.AdvanceChar()
		tok.Type = token.COMMA
	case '.':
		if IsDigit(s.PeekChar()) {
			tok.Type = token.NUMBER
			if err := scanNumber(s); err != nil {
				tok.Type = token.ILLEGAL
			}
		} else {
			s.AdvanceChar()
			tok.Type = token.DOT
		}
	case ';':
		s.AdvanceChar()
		tok.Type = token.SEMICOLON
	case '(':
		s.AdvanceChar()
		tok.Type = token.LPAREN
	case ')':
		s.AdvanceChar()
		tok.Type = token.RPAREN
	case '{':
		s.AdvanceChar()
		tok.Type = token.LBRACE
	case '}':
		s.AdvanceChar()
		tok.Type = token.RBRACE
	case '[':
		s.AdvanceChar()
		tok.Type = token.LBRACKET
	case ']':
		s.AdvanceChar()
		tok.Type = token.RBRACKET
	case ':':
		s.AdvanceChar()
		tok.Type = token.COLON
	case '\r':
		s.AdvanceChar()
		if s.currentChar == '\n' {
			s.AdvanceChar()
			tok.Type = token.NEWLINE
		} else {
			tok.Type = token.NEWLINE
		}
	case '\n':
		s.AdvanceChar()
		tok.Type = token.NEWLINE
	default:
		if IsLetter(s.currentChar) {
			scanIdentifier(s)
			tok.Type = token.IDENT
		} else if IsDigit(s.currentChar) {
			tok.Type = token.NUMBER
			if err := scanNumber(s); err != nil {
				tok.Type = token.ILLEGAL
			}
		} else if s.currentChar == utf8.RuneError {
			s.AdvanceChar()
			tok.Type = token.ILLEGAL
		} else if s.currentChar == EOF {
			tok.Type = token.EOF
		} else {
			s.AdvanceChar()
			tok.Type = token.UNKNOWN
		}
	}
	tok.Literal = s.input[ofs0 : s.offset-s.currentSize]
	return
}
