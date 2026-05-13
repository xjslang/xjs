package scanner

import (
	"strings"
	"unicode/utf8"
)

const eof = rune(-1)

type Scanner struct {
	input        []byte
	offset       int
	line, column int

	scanner func(*Scanner) Token

	CurrentChar rune
}

func (sc *Scanner) Init(input []byte) {
	sc.input = input
	if sc.scanner == nil {
		sc.scanner = defaultScanner
	}
	sc.Reset()
}

func (sc *Scanner) Reset() {
	if sc.scanner == nil {
		sc.scanner = defaultScanner
	}
	sc.offset = 0
	sc.CurrentChar = eof
	sc.line = 0
	sc.column = -1
	sc.AdvanceChar()
}

func (sc *Scanner) PeekChar() rune {
	if sc.offset < len(sc.input) {
		r, _ := utf8.DecodeRune(sc.input[sc.offset:])
		return r
	}
	return eof
}

func (sc *Scanner) AdvanceChar() {
	r, size := utf8.DecodeRune(sc.input[sc.offset:])
	sc.offset += size
	// covers "\r", "\n" and "\r\n"
	switch r {
	case '\r':
		sc.line++
		sc.column = -1
	case '\n':
		if sc.CurrentChar != '\r' {
			sc.line++
			sc.column = -1
		}
	case utf8.RuneError:
		if size > 0 {
			// just an illegal character; keep going
			sc.column++
		} else {
			// reached the end of the file
			r = eof
		}
	default:
		sc.column++
	}
	sc.CurrentChar = r
}

func (sc *Scanner) NextToken() Token {
	next := func() Token {
		sc.skipWhitespaces()
		line, column := sc.line, sc.column
		tok := sc.scanner(sc)
		tok.Line = line
		tok.Column = max(0, column)
		return tok
	}
	var trivia []Token
	afterNewline := false
	tok := next()
triviaLoop:
	for {
		switch tok.Type {
		case NEWLINE:
			afterNewline = true
		case LINE_COMMENT, BLOCK_COMMENT:
			afterNewline = afterNewline || strings.ContainsAny(tok.Literal, "\n\r")
		default:
			break triviaLoop
		}
		trivia = append(trivia, tok)
		tok = next()
	}
	tok.LeadingTrivia = trivia
	tok.AfterNewline = afterNewline
	return tok
}

func (sc *Scanner) skipWhitespaces() {
	for sc.CurrentChar == ' ' || sc.CurrentChar == '\t' {
		sc.AdvanceChar()
	}
}
