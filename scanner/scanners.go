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
		if sc.CurrentChar() == '\n' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			break
		} else if sc.currentChar == '\r' {
			sb.WriteRune(sc.currentChar)
			sc.AdvanceChar()
			if sc.CurrentChar() == '\n' {
				sb.WriteRune(sc.currentChar)
				sc.AdvanceChar()
			}
			break
		} else if sc.CurrentChar() == EOF {
			break
		}
		sb.WriteRune(sc.CurrentChar())
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
