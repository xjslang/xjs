package lexer

import (
	"io"

	"github.com/xjslang/xjs/token"
)

const eof = rune(-1)

type Lexer struct {
	input       io.RuneReader
	CurrentChar rune
	PeekChar    rune

	tokenReader func(l *Lexer) token.Token
}

func New(input io.RuneReader) *Lexer {
	l := &Lexer{
		input:       input,
		CurrentChar: eof,
		PeekChar:    eof,
		tokenReader: defaultTokenReader,
	}
	// called twice to init both CurrentChar and PeekChar
	l.ReadChar()
	l.ReadChar()
	return l
}

func (l *Lexer) ReadChar() {
	l.CurrentChar = l.PeekChar
	r, _, err := l.input.ReadRune()
	if err == io.EOF {
		l.PeekChar = eof
		return
	}
	if err != nil {
		panic(err)
	}
	l.PeekChar = r
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaces()
	return l.tokenReader(l)
}

func (l *Lexer) skipWhitespaces() {
	for isWhitespace(l.CurrentChar) {
		l.ReadChar()
	}
}
