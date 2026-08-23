package printer

import "github.com/xjslang/xjs/token"

func ErrorAt(offset int, msg string) error {
	return token.Error{
		Range:   token.Range{offset, offset},
		Message: msg,
	}
}

func isNewLine(r rune) bool {
	return r == eol || r == '\r' || r == '\n'
}

func isWhitespace(r rune) bool {
	return isNewLine(r) || r == ' ' || r == '\t'
}
