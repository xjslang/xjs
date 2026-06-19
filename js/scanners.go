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
