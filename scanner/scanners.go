package scanner

import (
	"errors"
	"strings"
)

func scanIdentifier(sc *Scanner) string {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); isLetter(sc.currentChar) || isDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String()
}

func scanHexNumber(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume 0
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume x | X
	if !isHexDigit(sc.currentChar) {
		return sb.String(), errors.New("hex digit expected")
	}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); isHexDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String(), nil
}

func scanOctalNumber(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume 0
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume o | O
	if !isOctalDigit(sc.currentChar) {
		return sb.String(), errors.New("octal digit expected")
	}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); isOctalDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String(), nil
}

func scanFloatNumber(sc *Scanner) (string, error) {
	sb := strings.Builder{}
	readDigits := func() {
		for sc.AdvanceChar(); isDigit(sc.currentChar); sc.AdvanceChar() {
			sb.WriteRune(sc.currentChar)
		}
	}
	sb.WriteRune(sc.currentChar)
	readDigits()
	if sc.currentChar == '.' {
		sb.WriteRune(sc.currentChar)
		readDigits()
	}
	if sc.currentChar == 'e' || sc.currentChar == 'E' {
		sb.WriteRune(sc.currentChar)
		sc.AdvanceChar()
		if sc.currentChar == '+' || sc.currentChar == '-' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
		}
		if !isDigit(sc.currentChar) {
			return sb.String(), errors.New("decimal digit expected")
		}
		sb.WriteRune(sc.currentChar)
		readDigits()
	}
	return sb.String(), nil
}

func scanNumber(sc *Scanner) (string, error) {
	if sc.currentChar == '0' {
		switch sc.PeekChar() {
		case 'x', 'X':
			return scanHexNumber(sc)
		case 'o', 'O':
			return scanOctalNumber(sc)
		}
	}
	return scanFloatNumber(sc)
}

func scanString(sc *Scanner, delimiter rune) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	for {
		sc.AdvanceChar()
		if sc.currentChar == '\\' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			if sc.currentChar == delimiter {
				sb.WriteRune(sc.currentChar)
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
	}
	return sb.String(), nil
}
