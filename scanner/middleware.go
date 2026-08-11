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
			tok = token.Token{Type: token.INCREMENT, Literal: "++"}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: "+"}
		}
	case '-':
		s.AdvanceChar()
		if s.currentChar == '-' {
			s.AdvanceChar()
			tok = token.Token{Type: token.DECREMENT, Literal: "--"}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: "-"}
		}
	case '*':
		s.AdvanceChar()
		tok = token.Token{Type: token.MULTIPLY, Literal: "*"}
	case '%':
		s.AdvanceChar()
		tok = token.Token{Type: token.MODULO, Literal: "%"}
	// divide operator and comments
	case '/':
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
			tok = token.Token{Type: token.DIVIDE, Literal: "/"}
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
		s.AdvanceChar()
		tok = token.Token{Type: token.COMMA, Literal: ","}
	case '.':
		if IsDigit(s.PeekChar()) {
			tok.Type = token.NUMBER
			if tok.Literal, err = ScanNumber(s); err != nil {
				tok.Type = token.ILLEGAL
				return
			}
		} else {
			s.AdvanceChar()
			tok = token.Token{Type: token.DOT, Literal: "."}
		}
	case ';':
		s.AdvanceChar()
		tok = token.Token{Type: token.SEMICOLON, Literal: ";"}
	case '(':
		s.AdvanceChar()
		tok = token.Token{Type: token.LPAREN, Literal: "("}
	case ')':
		s.AdvanceChar()
		tok = token.Token{Type: token.RPAREN, Literal: ")"}
	case '{':
		s.AdvanceChar()
		tok = token.Token{Type: token.LBRACE, Literal: "{"}
	case '}':
		s.AdvanceChar()
		tok = token.Token{Type: token.RBRACE, Literal: "}"}
	case '[':
		s.AdvanceChar()
		tok = token.Token{Type: token.LBRACKET, Literal: "["}
		return
	case ']':
		s.AdvanceChar()
		tok = token.Token{Type: token.RBRACKET, Literal: "]"}
	case ':':
		s.AdvanceChar()
		tok = token.Token{Type: token.COLON, Literal: ":"}
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
