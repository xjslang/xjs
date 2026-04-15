package lexer

import (
	"io"
	"strings"

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
