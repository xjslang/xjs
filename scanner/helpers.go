package scanner

import "unicode"

func IsLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || r == '$' || r > 127 && unicode.IsLetter(r)
}

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
