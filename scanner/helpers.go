package scanner

func IsLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '$'
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
