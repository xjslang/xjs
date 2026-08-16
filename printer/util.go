package printer

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
