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
			if sc.currentChar == delimiter {
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
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
	sb.WriteRune('`')
	for {
		if sc.currentChar == '`' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			break
		} else if sc.currentChar == EOF {
			return sb.String(), errors.New("unexpected end of file")
		}
		sb.WriteRune(sc.currentChar)
		sc.AdvanceChar()
	}
	return sb.String(), nil
}

func ScanHexNumber(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume x | X
	if !IsHexDigit(sc.currentChar) {
		return sb.String(), errors.New("hex digit expected")
	}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); IsHexDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String(), nil
}

func ScanOctalNumber(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume o | O
	if !IsOctalDigit(sc.currentChar) {
		return sb.String(), errors.New("octal digit expected")
	}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); IsOctalDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String(), nil
}

func ScanNumber(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	readDigits := func() {
		for sc.AdvanceChar(); IsDigit(sc.currentChar); sc.AdvanceChar() {
			sb.WriteRune(sc.currentChar)
		}
	}
	sb.WriteRune(sc.currentChar)
	readDigits()
	if sc.currentChar == '.' {
		sb.WriteRune(sc.currentChar)
		readDigits()
	}
	if c := sc.currentChar; c == 'e' || c == 'E' {
		sb.WriteRune(c)
		sc.AdvanceChar()
		if c := sc.currentChar; c == '+' || c == '-' {
			sb.WriteRune(c)
			sc.AdvanceChar()
		}
		if !IsDigit(sc.currentChar) {
			return sb.String(), errors.New("decimal digit expected")
		}
		sb.WriteRune(sc.currentChar)
		readDigits()
	}
	return sb.String(), nil
}
