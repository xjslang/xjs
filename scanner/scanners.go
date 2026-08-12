package scanner

import (
	"errors"
)

func ScanIdentifier(sc *Scanner) {
	for sc.AdvanceChar(); IsLetter(sc.currentChar) || IsDigit(sc.currentChar); sc.AdvanceChar() {
	}
}

func ScanLineComment(sc *Scanner) {
	for {
		sc.AdvanceChar()
		if sc.currentChar == '\n' {
			sc.AdvanceChar()
			break
		} else if sc.currentChar == '\r' {
			sc.AdvanceChar()
			if sc.currentChar == '\n' {
				sc.AdvanceChar()
			}
			break
		} else if sc.currentChar == EOF {
			break
		}
	}
}

func ScanBlockComment(sc *Scanner) error {
	sc.AdvanceChar() // consume "*"
	for {
		if sc.currentChar == '*' && sc.PeekChar() == '/' {
			// consume "*/"
			for range 2 {
				sc.AdvanceChar()
			}
			break
		} else if sc.currentChar == EOF {
			return errors.New("unexpected end of file")
		}
		sc.AdvanceChar()
	}
	return nil
}

func ScanString(sc *Scanner, delimiter rune) error {
	for {
		if sc.currentChar == '\\' {
			sc.AdvanceChar()
			if sc.currentChar == delimiter || sc.currentChar == '\n' {
				sc.AdvanceChar()
				continue
			} else if sc.currentChar == '\r' {
				sc.AdvanceChar()
				if sc.currentChar == '\n' {
					sc.AdvanceChar()
				}
				continue
			}
		}
		if sc.currentChar == delimiter {
			sc.AdvanceChar()
			break
		} else if sc.currentChar == EOF || sc.currentChar == '\n' || sc.currentChar == '\r' {
			return errors.New("unexpected end of line")
		}
		sc.AdvanceChar()
	}
	return nil
}

func ScanRawString(sc *Scanner) error {
	scanHole := func() error {
		depth := 1
		for {
			if sc.currentChar == EOF {
				return errors.New("unexpected end of file")
			} else if sc.currentChar == '`' {
				sc.AdvanceChar()
				err := ScanRawString(sc)
				if err != nil {
					return err
				}
				continue
			} else if sc.currentChar == '{' {
				depth++
			} else if sc.currentChar == '}' {
				depth--
				sc.AdvanceChar()
				if depth == 0 {
					return nil
				}
				continue
			}
			sc.AdvanceChar()
		}
	}

	for {
		if sc.currentChar == '`' {
			sc.AdvanceChar()
			break
		} else if sc.currentChar == '\\' {
			sc.AdvanceChar()
			switch sc.currentChar {
			case '`', '$':
				sc.AdvanceChar()
				continue
			}
		} else if sc.currentChar == '$' {
			sc.AdvanceChar()
			if sc.currentChar == '{' {
				sc.AdvanceChar()
				if err := scanHole(); err != nil {
					return err
				}
			}
			continue
		} else if sc.currentChar == EOF {
			return errors.New("unexpected end of file")
		}
		sc.AdvanceChar()
	}
	return nil
}

func ScanNumber(sc *Scanner) (err error) {
	scanDigits := func(check func(rune) bool) error {
		sc.AdvanceChar()
		if !check(sc.currentChar) {
			return errors.New("digit expected")
		}
		for check(sc.currentChar) {
			sc.AdvanceChar()
		}
		if IsLetter(sc.currentChar) || IsDigit(sc.currentChar) {
			return errors.New("digit expected")
		}
		return nil
	}
	scanDecimal := func() error {
		for IsDigit(sc.currentChar) {
			sc.AdvanceChar()
		}
		if sc.currentChar == '.' {
			sc.AdvanceChar()
		}
		for IsDigit(sc.currentChar) {
			sc.AdvanceChar()
		}
		if c := sc.currentChar; c == 'e' || c == 'E' {
			sc.AdvanceChar()
			if c := sc.currentChar; c == '+' || c == '-' {
				sc.AdvanceChar()
			}
			if !IsDigit(sc.currentChar) {
				return errors.New("digit expected")
			}
			for {
				sc.AdvanceChar()
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
		sc.AdvanceChar()
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
					sc.AdvanceChar()
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
			sc.AdvanceChar()
		}
		return err
	}
	return nil
}
