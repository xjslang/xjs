package lexer

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '$'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// hexDigitValue returns the numeric value of a hexadecimal digit
func hexDigitValue(ch byte) int {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0')
	}
	if ch >= 'a' && ch <= 'f' {
		return int(ch - 'a' + 10)
	}
	if ch >= 'A' && ch <= 'F' {
		return int(ch - 'A' + 10)
	}
	return 0
}

// encodeUTF8 converts a Unicode code point to UTF-8 byte sequence
func encodeUTF8(codePoint int) []byte {
	if codePoint <= 0x7F {
		// 1-byte sequence (ASCII)
		return []byte{byte(codePoint)}
	} else if codePoint <= 0x7FF {
		// 2-byte sequence
		return []byte{
			0xC0 | byte(codePoint>>6),
			0x80 | byte(codePoint&0x3F),
		}
	} else if codePoint <= 0xFFFF {
		// 3-byte sequence
		return []byte{
			0xE0 | byte(codePoint>>12),
			0x80 | byte((codePoint>>6)&0x3F),
			0x80 | byte(codePoint&0x3F),
		}
	} else if codePoint <= 0x10FFFF {
		// 4-byte sequence
		return []byte{
			0xF0 | byte(codePoint>>18),
			0x80 | byte((codePoint>>12)&0x3F),
			0x80 | byte((codePoint>>6)&0x3F),
			0x80 | byte(codePoint&0x3F),
		}
	}
	// Invalid code point, return replacement character (U+FFFD)
	return []byte{0xEF, 0xBF, 0xBD}
}
