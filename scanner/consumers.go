package scanner

import (
	"strings"

	"github.com/xjslang/xjs/token"
)

func (sc *Scanner) consumeBlockComment() (string, token.Type) {
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

func (sc *Scanner) consumeLineComment() string {
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

func (sc *Scanner) consumeIdentifier() string {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	for sc.AdvanceChar(); isLetter(sc.currentChar) || isDigit(sc.currentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.currentChar)
	}
	return sb.String()
}

func (sc *Scanner) consumeNumber() (string, token.Type) {
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

func (sc *Scanner) consumeString(delimiter rune) (string, token.Type) {
	sb := strings.Builder{}
	sb.WriteRune(sc.currentChar)
	for {
		sc.AdvanceChar()
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
