package js

import (
	"errors"
	"strings"

	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	LINE_COMMENT  = token.RegisterType("//")
	BLOCK_COMMENT = token.RegisterType("/*")
	NUMBER        = token.RegisterType("number")
	STRING        = token.RegisterType("string")
)

func ScanLineComment(s *scanner.Scanner) string {
	sb := strings.Builder{}
	sb.WriteString("//")
	for {
		s.AdvanceChar()
		if s.CurrentChar() == '\n' {
			sb.WriteRune(s.CurrentChar())
			s.AdvanceChar()
			break
		} else if s.CurrentChar() == '\r' {
			sb.WriteRune(s.CurrentChar())
			s.AdvanceChar()
			if s.CurrentChar() == '\n' {
				sb.WriteRune(s.CurrentChar())
				s.AdvanceChar()
			}
			break
		} else if s.CurrentChar() == scanner.EOF {
			break
		}
		sb.WriteRune(s.CurrentChar())
	}
	return sb.String()
}

func ScanBlockComment(sc *scanner.Scanner) (string, error) {
	sb := strings.Builder{}
	sb.WriteString("/*")
	sc.AdvanceChar() // consume "*"
	for {
		if sc.CurrentChar() == '*' && sc.PeekChar() == '/' {
			// consume "*/"
			for range 2 {
				sb.WriteRune(sc.CurrentChar())
				sc.AdvanceChar()
			}
			break
		} else if sc.CurrentChar() == scanner.EOF {
			return sb.String(), errors.New("unexpected end of file")
		}
		sb.WriteRune(sc.CurrentChar())
		sc.AdvanceChar()
	}
	return sb.String(), nil
}

func ScanString(sc *scanner.Scanner, delimiter rune) (string, error) {
	sb := strings.Builder{}
	for {
		if sc.CurrentChar() == '\\' {
			sb.WriteRune(sc.CurrentChar())
			sc.AdvanceChar()
			if sc.CurrentChar() == delimiter {
				sb.WriteRune(sc.CurrentChar())
				sc.AdvanceChar()
				continue
			}
		}
		if sc.CurrentChar() == delimiter {
			sb.WriteRune(sc.CurrentChar())
			sc.AdvanceChar()
			break
		} else if sc.CurrentChar() == scanner.EOF || sc.CurrentChar() == '\n' || sc.CurrentChar() == '\r' {
			return sb.String(), errors.New("unexpected end of line")
		}
		sb.WriteRune(sc.CurrentChar())
		sc.AdvanceChar()
	}
	return sb.String(), nil
}

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
