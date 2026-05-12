package scanner

import (
	"strings"
)

func (sc *Scanner) consumeBlockComment() (string, Kind) {
	sb := strings.Builder{}
	sc.AdvanceChar() // consume "*"
	for {
		if sc.CurrentChar == '*' && sc.PeekChar() == '/' {
			// skip "*/"
			sc.AdvanceChar()
			sc.AdvanceChar()
			break
		} else if sc.CurrentChar == eof {
			return sb.String(), ILLEGAL
		}
		sb.WriteRune(sc.CurrentChar)
		sc.AdvanceChar()
	}
	return sb.String(), BLOCK_COMMENT
}

func (sc *Scanner) consumeLineComment() string {
	sb := strings.Builder{}
	for {
		sc.AdvanceChar()
		if sc.CurrentChar == '\n' {
			sb.WriteRune(sc.CurrentChar)
			sc.AdvanceChar()
			break
		} else if sc.CurrentChar == '\r' {
			sb.WriteRune(sc.CurrentChar)
			sc.AdvanceChar()
			if sc.CurrentChar == '\n' {
				sb.WriteRune(sc.CurrentChar)
				sc.AdvanceChar()
			}
			break
		} else if sc.CurrentChar == eof {
			break
		}
		sb.WriteRune(sc.CurrentChar)
	}
	return sb.String()
}

func (sc *Scanner) consumeIdentifier() string {
	sb := strings.Builder{}
	sb.WriteRune(sc.CurrentChar)
	for sc.AdvanceChar(); isLetter(sc.CurrentChar) || isDigit(sc.CurrentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.CurrentChar)
	}
	return sb.String()
}

func (sc *Scanner) consumeNumber() string {
	sb := strings.Builder{}
	sb.WriteRune(sc.CurrentChar)
	for sc.AdvanceChar(); isDigit(sc.CurrentChar); sc.AdvanceChar() {
		sb.WriteRune(sc.CurrentChar)
	}
	return sb.String()
}

func (sc *Scanner) consumeString(delimiter rune) (string, Kind) {
	sb := strings.Builder{}
	sb.WriteRune(sc.CurrentChar)
	for {
		sc.AdvanceChar()
		if sc.CurrentChar == delimiter {
			sb.WriteRune(sc.CurrentChar)
			sc.AdvanceChar()
			break
		} else if sc.CurrentChar == eof || sc.CurrentChar == '\n' || sc.CurrentChar == '\r' {
			return sb.String(), ILLEGAL
		}
		sb.WriteRune(sc.CurrentChar)
	}
	return sb.String(), STRING
}
