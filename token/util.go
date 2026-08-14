package token

import "unicode/utf8"

func Position(src string, offset int) (line, col int) {
	if offset < 0 {
		return 0, 0
	}
	if offset > len(src) {
		offset = len(src)
	}
	i := 0
	var prevR rune
	for i < offset {
		r, size := utf8.DecodeRuneInString(src[i:])
		if r == utf8.RuneError && size == 0 {
			break
		}
		if i+size > offset {
			break
		}
		i += size
		switch r {
		case '\n':
			if prevR != '\r' {
				line++
				col = 0
			}
		case '\r':
			line++
			col = 0
		default:
			col++
		}
		prevR = r
	}
	return
}
