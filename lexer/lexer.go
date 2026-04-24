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

	tokenizer func(l *Lexer) token.Token

	CurrentChar rune
}

func (l *Lexer) Init(input []byte) {
	l.input = input
	l.tokenizer = defaultTokenizer
	l.Reset()
}

func (l *Lexer) Reset() {
	if l.tokenizer == nil {
		l.tokenizer = defaultTokenizer
	}
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
	r, size := utf8.DecodeRune(l.input[l.offset:])
	l.offset += size
	// covers "\r", "\n" and "\r\n"
	switch r {
	case '\r':
		l.line++
		l.column = -1
	case '\n':
		if l.CurrentChar != '\r' {
			l.line++
			l.column = -1
		}
	case utf8.RuneError:
		if size > 0 {
			// just an illegal character; keep going
			l.column++
		} else {
			// reached the end of the file
			r = eof
		}
	default:
		l.column++
	}
	l.CurrentChar = r
}

func (l *Lexer) NextToken() token.Token {
	next := func() token.Token {
		l.skipWhitespaces()
		line, column := l.line, l.column
		tok := l.tokenizer(l)
		tok.Line = line
		tok.Column = max(0, column)
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
