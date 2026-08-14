package printer

import (
	"strings"
)

func ErrorAt(offset int, msg string) error {
	return Error{
		Offset:  offset,
		Message: msg,
	}
}

func isNewLine(r rune) bool {
	return r == eol || r == '\r' || r == '\n'
}

func isWhitespace(r rune) bool {
	return isNewLine(r) || r == ' ' || r == '\t'
}

func splitTrivia(trivia string) []string {
	var lines []string
	sb := strings.Builder{}
	skipSpaces := true
	skipLF := false
	for _, c := range trivia {
		if skipLF {
			skipLF = false
			if c == '\n' {
				continue
			}
		}
		if skipSpaces && (c == ' ' || c == '\t') {
			continue
		}
		switch c {
		case '\r':
			sb.WriteRune('\n') // normalize CR/CRLF to LF
			lines = append(lines, sb.String())
			sb.Reset()
			skipSpaces = true
			skipLF = true
			continue
		case '\n':
			sb.WriteRune('\n')
			lines = append(lines, sb.String())
			sb.Reset()
			skipSpaces = true
			continue
		default:
			sb.WriteRune(c)
			skipSpaces = false
		}
	}
	if s := sb.String(); len(s) > 0 {
		lines = append(lines, strings.TrimRight(s, " "))
	}
	return lines
}
