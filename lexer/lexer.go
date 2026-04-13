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
	// call twice to update CurrentChar and PeekChar
	l.AdvanceChar()
	l.AdvanceChar()
	return l
}

func (l *Lexer) AdvanceChar() {
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
	next := func() token.Token {
		l.skipWhitespaces()
		return l.tokenReader(l)
	}
	var trivia []string
	tok := next()
	for ; tok.Type == token.NEWLINE || tok.Type == token.SINGLELINE_COMMENT || tok.Type == token.MULTILINE_COMMENT; tok = next() {
		trivia = append(trivia, tok.Literal)
	}
	tok.LeadingTrivia = trivia
	return tok
}

func (l *Lexer) skipWhitespaces() {
	for l.CurrentChar == ' ' || l.CurrentChar == '\t' || l.CurrentChar == '\r' {
		l.AdvanceChar()
	}
}
