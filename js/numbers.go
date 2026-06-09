package js

import (
	"strings"

	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func ScanHexNumber(sc *scanner.Scanner) (string, token.Type) {
	if c := sc.CurrentChar(); c != 'x' && c != 'X' {
		return "", token.ILLEGAL
	}
	sb := strings.Builder{}
	sb.WriteRune(sc.CurrentChar())
	sc.AdvanceChar() // consume x | X
	if !isHexDigit(sc.CurrentChar()) {
		return sb.String(), token.ILLEGAL
	}
	sb.WriteRune(sc.CurrentChar())
	for sc.AdvanceChar(); isHexDigit(sc.CurrentChar()); sc.AdvanceChar() {
		sb.WriteRune(sc.CurrentChar())
	}
	return sb.String(), token.NUMBER
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexDigit(r rune) bool {
	return isDigit(r) || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}
