package lexer

import (
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

const eof = rune(-1)

type Lexer struct {
	input        []byte
	offset       int
	line, column int

	tokenReader func(l *Lexer) token.Token

	CurrentChar rune
}

func New(input []byte) *Lexer {
	l := &Lexer{
		input:       input,
		tokenReader: defaultTokenReader,
	}
	l.Reset()
	return l
}

func (l *Lexer) Reset() {
	l.offset = 0
	l.CurrentChar = eof
	l.line = 0
	l.column = -1
	l.AdvanceChar()
}

func (l *Lexer) PeekChar() rune {
	if l.offset < len(l.input) {
		r, _ := utf8.DecodeRune(l.input[l.offset:])
		return r
	}
	return eof
}

func (l *Lexer) AdvanceChar() {
	if l.offset < len(l.input) {
		r, size := utf8.DecodeRune(l.input[l.offset:])
		l.offset += size
		// the following condition covers "\r", "\n" and "\r\n"
		if r == '\r' || (l.CurrentChar != '\r' && r == '\n') {
			l.line++
			l.column = -1
		} else if r != '\n' {
			l.column++
		}
		l.CurrentChar = r
	} else {
		if l.column < 0 {
			l.column = 0
		}
		l.CurrentChar = eof
	}
}

func (l *Lexer) NextToken() token.Token {
	next := func() token.Token {
		l.skipWhitespaces()
		line, column := l.line, l.column
		tok := l.tokenReader(l)
		tok.Line = line
		tok.Column = column
		return tok
	}
	var trivia []string
	afterNewline := false
	tok := next()
triviaLoop:
	for {
		switch tok.Type {
		case token.NEWLINE:
			afterNewline = true
		case token.LINE_COMMENT:
		case token.BLOCK_COMMENT:
			afterNewline = afterNewline || strings.ContainsAny(tok.Literal, "\n\r")
		default:
			break triviaLoop
		}
		trivia = append(trivia, tok.Literal)
		tok = next()
	}
	tok.LeadingTrivia = trivia
	tok.AfterNewline = afterNewline
	return tok
}

func (l *Lexer) skipWhitespaces() {
	for l.CurrentChar == ' ' || l.CurrentChar == '\t' {
		l.AdvanceChar()
	}
}
