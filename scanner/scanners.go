package scanner

import (
	"github.com/xjslang/xjs/token"
)

type Error struct {
	token.Type
}

func (e Error) Error() string {
	return e.Type.Name()
}

func scanIdentifier(s *Scanner) {
	for IsLetter(s.currentChar) || IsDigit(s.currentChar) {
		s.AdvanceChar()
	}
}

func scanLineComment(s *Scanner) {
	for {
		s.AdvanceChar()
		if s.currentChar == '\n' {
			s.AdvanceChar()
			break
		} else if s.currentChar == '\r' {
			s.AdvanceChar()
			if s.currentChar == '\n' {
				s.AdvanceChar()
			}
			break
		} else if s.currentChar == EOF {
			break
		}
	}
}

func scanBlockComment(s *Scanner) error {
	s.AdvanceChar() // consume "*"
	for {
		if s.currentChar == '*' && s.PeekChar() == '/' {
			// consume "*/"
			for range 2 {
				s.AdvanceChar()
			}
			break
		} else if s.currentChar == EOF {
			return Error{Type: token.UNCLOSED_BLOCK_COMMENT}
		}
		s.AdvanceChar()
	}
	return nil
}

func scanTrivia(s *Scanner) error {
loop:
	for {
		// skip spaces
		for {
			if c := s.currentChar; c != ' ' && c != '\t' && c != '\n' && c != '\r' {
				break
			}
			s.AdvanceChar()
		}
		if s.currentChar != '/' {
			break
		}
		switch s.PeekChar() {
		case '/':
			s.AdvanceChar()
			scanLineComment(s)
		case '*':
			s.AdvanceChar()
			if err := scanBlockComment(s); err != nil {
				return err
			}
		default:
			break loop
		}
	}
	return nil
}

func scanString(s *Scanner, delimiter rune) error {
	for {
		if s.currentChar == '\\' {
			s.AdvanceChar()
			if s.currentChar == delimiter || s.currentChar == '\n' {
				s.AdvanceChar()
				continue
			} else if s.currentChar == '\r' {
				s.AdvanceChar()
				if s.currentChar == '\n' {
					s.AdvanceChar()
				}
				continue
			}
		}
		if s.currentChar == delimiter {
			s.AdvanceChar()
			break
		} else if s.currentChar == EOF || s.currentChar == '\n' || s.currentChar == '\r' {
			return Error{Type: token.UNCLOSED_STRING}
		}
		s.AdvanceChar()
	}
	return nil
}

func scanRawString(s *Scanner) error {
	scanHole := func() error {
		depth := 1
		for {
			if s.currentChar == EOF {
				return Error{Type: token.UNCLOSED_STRING}
			} else if s.currentChar == '`' {
				s.AdvanceChar()
				err := scanRawString(s)
				if err != nil {
					return err
				}
				continue
			} else if s.currentChar == '{' {
				depth++
			} else if s.currentChar == '}' {
				depth--
				s.AdvanceChar()
				if depth == 0 {
					return nil
				}
				continue
			}
			s.AdvanceChar()
		}
	}

	for {
		if s.currentChar == '`' {
			s.AdvanceChar()
			break
		} else if s.currentChar == '\\' {
			s.AdvanceChar()
			switch s.currentChar {
			case '`', '$':
				s.AdvanceChar()
				continue
			}
		} else if s.currentChar == '$' {
			s.AdvanceChar()
			if s.currentChar == '{' {
				s.AdvanceChar()
				if err := scanHole(); err != nil {
					return err
				}
			}
			continue
		} else if s.currentChar == EOF {
			return Error{Type: token.UNCLOSED_STRING}
		}
		s.AdvanceChar()
	}
	return nil
}

func scanNumber(s *Scanner) (err error) {
	scanDigits := func(check func(rune) bool) error {
		s.AdvanceChar()
		if !check(s.currentChar) {
			return Error{Type: token.INVALID_NUMBER}
		}
		for check(s.currentChar) {
			s.AdvanceChar()
		}
		if IsLetter(s.currentChar) || IsDigit(s.currentChar) {
			return Error{Type: token.INVALID_NUMBER}
		}
		return nil
	}
	scanDecimal := func() error {
		for IsDigit(s.currentChar) {
			s.AdvanceChar()
		}
		if s.currentChar == '.' {
			s.AdvanceChar()
		}
		for IsDigit(s.currentChar) {
			s.AdvanceChar()
		}
		if c := s.currentChar; c == 'e' || c == 'E' {
			s.AdvanceChar()
			if c := s.currentChar; c == '+' || c == '-' {
				s.AdvanceChar()
			}
			if !IsDigit(s.currentChar) {
				return Error{Type: token.INVALID_NUMBER}
			}
			for {
				s.AdvanceChar()
				if !IsDigit(s.currentChar) {
					break
				}
			}
		}
		if IsLetter(s.currentChar) {
			return Error{Type: token.INVALID_NUMBER}
		}
		return nil
	}
	switch s.currentChar {
	case '0':
		s.AdvanceChar()
		switch s.currentChar {
		case 'x', 'X':
			err = scanDigits(IsHexDigit)
		case 'o', 'O':
			err = scanDigits(IsOctalDigit)
		case 'b', 'B':
			err = scanDigits(IsBinaryDigit)
		default:
			if IsDigit(s.currentChar) {
				for IsDigit(s.currentChar) {
					s.AdvanceChar()
				}
				if c := s.currentChar; (c == '.' && IsDigit(s.PeekChar())) || c == 'e' || c == 'E' {
					err = scanDecimal()
				}
			} else {
				err = scanDecimal()
			}
		}
	default:
		err = scanDecimal()
	}
	if err != nil {
		for IsLetter(s.currentChar) || IsDigit(s.currentChar) {
			s.AdvanceChar()
		}
		return err
	}
	return nil
}
