package lexer

import "strings"

func (l *Lexer) consumeChar(sb *strings.Builder) {
	sb.WriteRune(l.CurrentChar)
	l.Advance()
}

func (l *Lexer) consumeChars(sb *strings.Builder, count int) {
	for range count {
		sb.WriteRune(l.CurrentChar)
		l.Advance()
	}
}

func (l *Lexer) consumeIdentifier(sb *strings.Builder) {
	sb.WriteRune(l.CurrentChar)
	for l.Advance(); isLetter(l.CurrentChar) || isDigit(l.CurrentChar); l.Advance() {
		sb.WriteRune(l.CurrentChar)
	}
}

func (l *Lexer) consumeNumber(sb *strings.Builder) {
	sb.WriteRune(l.CurrentChar)
	for l.Advance(); isDigit(l.CurrentChar); l.Advance() {
		sb.WriteRune(l.CurrentChar)
	}
}
