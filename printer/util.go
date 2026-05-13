package printer

func isNewLine(r rune) bool {
	return r == eol || r == '\r' || r == '\n'
}
