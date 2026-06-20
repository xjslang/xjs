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
