package lexer

import (
	"strings"

	"github.com/xjslang/xjs/token"
)

func (l *Lexer) consumeChar() string {
	lit := string(l.CurrentChar)
	l.ReadChar()
	return lit
}

func (l *Lexer) consumeChars(count int) string {
	sb := strings.Builder{}
	for range count {
		sb.WriteRune(l.CurrentChar)
		l.ReadChar()
	}
	return sb.String()
}

func (l *Lexer) consumeIdentifier() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.ReadChar(); isLetter(l.CurrentChar) || isDigit(l.CurrentChar); l.ReadChar() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}

func (l *Lexer) consumeNumber() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.ReadChar(); isDigit(l.CurrentChar); l.ReadChar() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}

func (l *Lexer) consumeString(delimiter rune) (string, token.TokenType) {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for {
		l.ReadChar()
		if l.CurrentChar == delimiter {
			sb.WriteRune(l.CurrentChar)
			l.ReadChar()
			break
		} else if l.CurrentChar == eof || l.CurrentChar == '\n' {
			return sb.String(), token.ILLEGAL
		}
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String(), token.STRING
}
