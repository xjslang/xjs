package scanner

import (
	"errors"
	"strings"
)

func ScanIdentifier(sc *Scanner) string {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); IsLetter(sc.currentChar) || IsDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String()
}

func ScanLineComment(sc *Scanner) string {
	sb := strings.Builder{}
	sb.WriteString("//")
	for {
		sc.AdvanceChar()
		if sc.currentChar == '\n' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			break
		} else if sc.currentChar == '\r' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			if sc.currentChar == '\n' {
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
			}
			break
		} else if sc.currentChar == EOF {
			break
		}
		sb.WriteRune(sc.currentChar)
	}
	return sb.String()
}

func ScanBlockComment(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("/*")
	sc.AdvanceChar() // consume "*"
	for {
		if sc.currentChar == '*' && sc.PeekChar() == '/' {
			// consume "*/"
			for range 2 {
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
			}
			break
		} else if sc.currentChar == EOF {
			return sb.String(), errors.New("unexpected end of file")
		}
		sb.WriteRune(sc.currentChar)
		sc.AdvanceChar()
	}
	return sb.String(), nil
}

func ScanString(sc *Scanner, delimiter rune) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(delimiter)
	for {
		if sc.currentChar == '\\' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			if sc.currentChar == delimiter || sc.currentChar == '\n' {
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
				continue
			} else if sc.currentChar == '\r' {
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
				if sc.currentChar == '\n' {
					sb.WriteRune(sc.currentChar)
					sc.AdvanceChar()
				}
				continue
			}
		}
		if sc.currentChar == delimiter {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			break
		} else if sc.currentChar == EOF || sc.currentChar == '\n' || sc.currentChar == '\r' {
			return sb.String(), errors.New("unexpected end of line")
		}
		sb.WriteRune(sc.currentChar)
		sc.AdvanceChar()
	}
	return sb.String(), nil
}

func ScanRawString(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	scanHole := func() error {
		depth := 1
		for {
			if sc.currentChar == EOF {
				return errors.New("unexpected end of file")
			} else if sc.currentChar == '`' {
				sc.AdvanceChar()
				s, err := ScanRawString(sc)
				sb.WriteString(s)
				if err != nil {
					return err
				}
				continue
			} else if sc.currentChar == '{' {
				depth++
			} else if sc.currentChar == '}' {
				depth--
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
				if depth == 0 {
					return nil
				}
				continue
			}
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
		}
	}

	sb.WriteRune('`')
	for {
		if sc.currentChar == '`' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			break
		} else if sc.currentChar == '\\' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			switch sc.currentChar {
			case '`', '$':
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
				continue
			}
		} else if sc.currentChar == '$' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			if sc.currentChar == '{' {
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
				if err := scanHole(); err != nil {
					return sb.String(), err
				}
			}
			continue
		} else if sc.currentChar == EOF {
			return sb.String(), errors.New("unexpected end of file")
		}
		sb.WriteRune(sc.currentChar)
		sc.AdvanceChar()
	}
	return sb.String(), nil
}

func ScanNumber(sc *Scanner) (_ string, err error) {
	var sb strings.Builder
	next := func() {
		sb.WriteRune(sc.currentChar)
		sc.AdvanceChar()
	}
	scanDigits := func(check func(rune) bool) error {
		next()
		if !check(sc.currentChar) {
			return errors.New("digit expected")
		}
		for check(sc.currentChar) {
			next()
		}
		if IsLetter(sc.currentChar) || IsDigit(sc.currentChar) {
			return errors.New("digit expected")
		}
		return nil
	}
	scanDecimal := func() error {
		for IsDigit(sc.currentChar) {
			next()
		}
		if sc.currentChar == '.' {
			next()
		}
		for IsDigit(sc.currentChar) {
			next()
		}
		if c := sc.currentChar; c == 'e' || c == 'E' {
			next()
			if c := sc.currentChar; c == '+' || c == '-' {
				next()
			}
			if !IsDigit(sc.currentChar) {
				return errors.New("digit expected")
			}
			for {
				next()
				if !IsDigit(sc.currentChar) {
					break
				}
			}
		}
		if IsLetter(sc.currentChar) {
			return errors.New("digit expected")
		}
		return nil
	}
	switch sc.currentChar {
	case '0':
		next()
		switch sc.currentChar {
		case 'x', 'X':
			err = scanDigits(IsHexDigit)
		case 'o', 'O':
			err = scanDigits(IsOctalDigit)
		case 'b', 'B':
			err = scanDigits(IsBinaryDigit)
		default:
			if IsDigit(sc.currentChar) {
				for IsDigit(sc.currentChar) {
					next()
				}
				if c := sc.currentChar; (c == '.' && IsDigit(sc.PeekChar())) || c == 'e' || c == 'E' {
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
		for IsLetter(sc.currentChar) || IsDigit(sc.currentChar) {
			next()
		}
		return sb.String(), err
	}
	return sb.String(), nil
}
