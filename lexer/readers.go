package lexer

import "strings"

func (l *Lexer) readChar() string {
	lit := string(l.CurrentChar)
	l.Advance()
	return lit
}

func (l *Lexer) readChars(count int) string {
	sb := strings.Builder{}
	for range count {
		sb.WriteRune(l.CurrentChar)
		l.Advance()
	}
	return sb.String()
}

func (l *Lexer) readIden() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.Advance(); isLetter(l.CurrentChar) || isDigit(l.CurrentChar); l.Advance() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}

func (l *Lexer) readNumber() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.Advance(); isDigit(l.CurrentChar); l.Advance() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}
