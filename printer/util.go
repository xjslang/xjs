package printer

import "github.com/xjslang/xjs/printer/internal/formatter"

func isNewLine(r rune) bool {
	return r == formatter.EOL || r == '\r' || r == '\n'
}

func isWhitespace(r rune) bool {
	return isNewLine(r) || r == ' ' || r == '\t'
}
