package lexer

import (
	"strings"

	"github.com/xjslang/xjs/token"
)

func (l *Lexer) parseBlockComment() (string, token.TokenType) {
	sb := strings.Builder{}
	l.AdvanceChar() // consume "*"
	for {
		if l.CurrentChar == '*' && l.PeekChar == '/' {
			// skip "*/"
			l.AdvanceChar()
			l.AdvanceChar()
			break
		} else if l.CurrentChar == eof {
			return sb.String(), token.ILLEGAL
		}
		sb.WriteRune(l.CurrentChar)
		l.AdvanceChar()
	}
	return sb.String(), token.BLOCK_COMMENT
}

func (l *Lexer) parseLineComment() string {
	sb := strings.Builder{}
	for l.AdvanceChar(); l.CurrentChar != '\n' && l.CurrentChar != '\r' && l.CurrentChar != eof; l.AdvanceChar() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}

func (l *Lexer) parseIdentifier() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.AdvanceChar(); isLetter(l.CurrentChar) || isDigit(l.CurrentChar); l.AdvanceChar() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}

func (l *Lexer) parseNumber() string {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for l.AdvanceChar(); isDigit(l.CurrentChar); l.AdvanceChar() {
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String()
}

func (l *Lexer) parseString(delimiter rune) (string, token.TokenType) {
	sb := strings.Builder{}
	sb.WriteRune(l.CurrentChar)
	for {
		l.AdvanceChar()
		if l.CurrentChar == delimiter {
			sb.WriteRune(l.CurrentChar)
			l.AdvanceChar()
			break
		} else if l.CurrentChar == eof || l.CurrentChar == '\n' || l.CurrentChar == '\r' {
			return sb.String(), token.ILLEGAL
		}
		sb.WriteRune(l.CurrentChar)
	}
	return sb.String(), token.STRING
}
