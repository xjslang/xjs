package printer

import "github.com/xjslang/xjs/token"

func ErrorAt(tok token.Token, msg string) error {
	return &Error{
		Position: tok.Position,
		Message:  msg,
	}
}

func isNewLine(r rune) bool {
	return r == eol || r == '\r' || r == '\n'
}

func isWhitespace(r rune) bool {
	return isNewLine(r) || r == ' ' || r == '\t'
}
