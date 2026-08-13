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
	ofs0 := s.offset - s.currentSize
	switch s.currentChar {
	// operators
	case '=':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.EQ, "=="
		} else {
			tok.Type, tok.Literal = token.ASSIGN, "="
		}
	case '!':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.NOT_EQ, "!="
		} else {
			tok.Type, tok.Literal = token.NOT, "!"
		}
	case '<':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.LTE, "<="
		} else {
			tok.Type, tok.Literal = token.LT, "<"
		}
	case '>':
		s.AdvanceChar()
		if s.currentChar == '=' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.GTE, ">="
		} else {
			tok.Type, tok.Literal = token.GT, ">"
		}
	case '|':
		c := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '|' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.OR, "||"
		} else {
			tok.Type, tok.Literal = token.UNKNOWN, string(c)
		}
	case '&':
		c := s.currentChar
		s.AdvanceChar()
		if s.currentChar == '&' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.AND, "&&"
		} else {
			tok.Type, tok.Literal = token.UNKNOWN, string(c)
		}
	// maths operators
	case '+':
		s.AdvanceChar()
		if s.currentChar == '+' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.INCREMENT, "++"
		} else {
			tok.Type, tok.Literal = token.PLUS, "+"
		}
	case '-':
		s.AdvanceChar()
		if s.currentChar == '-' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.DECREMENT, "--"
		} else {
			tok.Type, tok.Literal = token.MINUS, "-"
		}
	case '*':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.MULTIPLY, "*"
	case '%':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.MODULO, "%"
	// divide operator and comments
	case '/':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.DIVIDE, "/"
	// delimiters
	case '\'', '"':
		c := s.currentChar
		s.AdvanceChar()
		tok.Type = token.STRING
		if err = scanString(s, c); err != nil {
			tok.Type = token.ILLEGAL
		}
		tok.Literal = string(s.input[ofs0 : s.offset-s.currentSize])
	case '`':
		s.AdvanceChar()
		tok.Type = token.STRING
		if err = scanRawString(s); err != nil {
			tok.Type = token.ILLEGAL
		}
		tok.Literal = string(s.input[ofs0 : s.offset-s.currentSize])
	case ',':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.COMMA, ","
	case '.':
		if IsDigit(s.PeekChar()) {
			tok.Type = token.NUMBER
			if err = scanNumber(s); err != nil {
				tok.Type = token.ILLEGAL
			}
			tok.Literal = string(s.input[ofs0 : s.offset-s.currentSize])
		} else {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.DOT, "."
		}
	case ';':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.SEMICOLON, ";"
	case '(':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.LPAREN, "("
	case ')':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.RPAREN, ")"
	case '{':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.LBRACE, "{"
	case '}':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.RBRACE, "}"
	case '[':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.LBRACKET, "["
		return
	case ']':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.RBRACKET, "]"
	case ':':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.COLON, ":"
	case '\r':
		s.AdvanceChar()
		if s.currentChar == '\n' {
			s.AdvanceChar()
			tok.Type, tok.Literal = token.NEWLINE, "\r\n"
		} else {
			tok.Type, tok.Literal = token.NEWLINE, "\r"
		}
	case '\n':
		s.AdvanceChar()
		tok.Type, tok.Literal = token.NEWLINE, "\n"
	default:
		if IsLetter(s.currentChar) {
			scanIdentifier(s)
			lit := string(s.input[ofs0 : s.offset-s.currentSize])
			tok.Type, tok.Literal = token.IDENT, lit
		} else if IsDigit(s.currentChar) {
			tok.Type = token.NUMBER
			if err = scanNumber(s); err != nil {
				tok.Type = token.ILLEGAL
			}
			tok.Literal = string(s.input[ofs0 : s.offset-s.currentSize])
		} else if s.currentChar == utf8.RuneError {
			c := s.currentChar
			s.AdvanceChar()
			tok.Type, tok.Literal = token.ILLEGAL, string(c)
		} else if s.currentChar == EOF {
			tok.Type = token.EOF
		} else {
			c := s.currentChar
			s.AdvanceChar()
			tok.Type, tok.Literal = token.UNKNOWN, string(c)
		}
	}
	return
}
