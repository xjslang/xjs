package scanner

import (
	"strings"

	"github.com/xjslang/xjs/token"
)

func scanBlockComment(sc *Scanner) (string, token.Type) {
	sb := strings.Builder{}
	sc.AdvanceChar() // consume "*"
	for {
		if sc.currentChar == '*' && sc.PeekChar() == '/' {
			// skip "*/"
			sc.AdvanceChar()
			sc.AdvanceChar()
			break
		} else if sc.currentChar == eof {
			return sb.String(), token.ILLEGAL
		}
		sb.WriteRune(sc.currentChar)
		sc.AdvanceChar()
	}
	return sb.String(), token.BLOCK_COMMENT
}

func scanLineComment(sc *Scanner) string {
	sb := strings.Builder{}
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
		} else if sc.currentChar == eof {
			break
		}
		sb.WriteRune(sc.currentChar)
	}
	return sb.String()
}

func scanIdentifier(sc *Scanner) string {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); isLetter(sc.currentChar) || isDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String()
}

func scanHexNumber(sc *Scanner) (string, token.Type) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume 0
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume x | X
	if !isHexDigit(sc.currentChar) {
		return sb.String(), token.ILLEGAL
	}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); isHexDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String(), token.NUMBER
}

func scanOctalNumber(sc *Scanner) (string, token.Type) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume 0
	sb.WriteRune(sc.currentChar)
	sc.AdvanceChar() // consume x | X
	if !isOctalDigit(sc.currentChar) {
		return sb.String(), token.ILLEGAL
	}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); isOctalDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String(), token.NUMBER
}

func scanFloatNumber(sc *Scanner) (string, token.Type) {
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
			return sb.String(), token.ILLEGAL
		}
		sb.WriteRune(sc.currentChar)
		readDigits()
	}
	return sb.String(), token.NUMBER
}

func scanNumber(sc *Scanner) (string, token.Type) {
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

func scanString(sc *Scanner, delimiter rune) (string, token.Type) {
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
		} else if sc.currentChar == eof || sc.currentChar == '\n' || sc.currentChar == '\r' {
			return sb.String(), token.ILLEGAL
		}
		sb.WriteRune(sc.currentChar)
	}
	return sb.String(), token.STRING
}
