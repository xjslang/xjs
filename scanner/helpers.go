package scanner

func IsLetter(r rune) bool {
	return r >= 'a' && r <= 'z' ||
		r >= 'A' && r <= 'Z' ||
		r == '_' || r == '$' ||
		r == '\u200c' || // ZWNJ
		r == '\u200d' // ZWJ
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func IsHexDigit(r rune) bool {
	return IsDigit(r) || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

func IsOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}
