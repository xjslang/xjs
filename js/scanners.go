package js

import (
	"errors"
	"strings"

	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var NUMBER = token.RegisterType("number")

func ScanHexNumber(sc *scanner.Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(sc.CurrentChar())
	sc.AdvanceChar() // consume x | X
	if !isHexDigit(sc.CurrentChar()) {
		return sb.String(), errors.New("hex digit expected")
	}
	sb.WriteRune(sc.CurrentChar())
	for sc.AdvanceChar(); isHexDigit(sc.CurrentChar()); sc.AdvanceChar() {
		sb.WriteRune(sc.CurrentChar())
	}
	return sb.String(), nil
}

func ScanOctalNumber(sc *scanner.Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteRune(sc.CurrentChar())
	sc.AdvanceChar() // consume o | O
	if !isOctalDigit(sc.CurrentChar()) {
		return sb.String(), errors.New("octal digit expected")
	}
	sb.WriteRune(sc.CurrentChar())
	for sc.AdvanceChar(); isOctalDigit(sc.CurrentChar()); sc.AdvanceChar() {
		sb.WriteRune(sc.CurrentChar())
	}
	return sb.String(), nil
}

func ScanNumber(sc *scanner.Scanner) (string, error) {
	sb := strings.Builder{}
	readDigits := func() {
		for sc.AdvanceChar(); isDigit(sc.CurrentChar()); sc.AdvanceChar() {
			sb.WriteRune(sc.CurrentChar())
		}
	}
	sb.WriteRune(sc.CurrentChar())
	readDigits()
	if sc.CurrentChar() == '.' {
		sb.WriteRune(sc.CurrentChar())
		readDigits()
	}
	if c := sc.CurrentChar(); c == 'e' || c == 'E' {
		sb.WriteRune(c)
		sc.AdvanceChar()
		if c := sc.CurrentChar(); c == '+' || c == '-' {
			sb.WriteRune(c)
			sc.AdvanceChar()
		}
		if !isDigit(sc.CurrentChar()) {
			return sb.String(), errors.New("decimal digit expected")
		}
		sb.WriteRune(sc.CurrentChar())
		readDigits()
	}
	return sb.String(), nil
}

func isHexDigit(r rune) bool {
	return isDigit(r) || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
