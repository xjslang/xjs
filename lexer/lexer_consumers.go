package lexer

import "strings"

func (l *Lexer) consumeChar() string {
	lit := string(l.CurrentChar)
	l.Advance()
	return lit
}

func (l *Lexer) consumeChars(count int) string {
	sb := strings.Builder{}
	for range count {
		sb.WriteRune(l.CurrentChar)
		l.Advance()
	}
	return sb.String()
}

func (l *Lexer) consumeIdentifier() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.Advance(); isLetter(l.CurrentChar) || isDigit(l.CurrentChar); l.Advance() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}

func (l *Lexer) consumeNumber() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.Advance(); isDigit(l.CurrentChar); l.Advance() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}
