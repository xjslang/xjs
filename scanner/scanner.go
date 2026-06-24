package scanner

import (
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

type config struct {
	withTriviaTypes []token.Type
}

func WithCommentTypes(typ ...token.Type) func(*config) {
	return func(cfg *config) {
		cfg.withTriviaTypes = append(cfg.withTriviaTypes, typ...)
	}
}

const EOF = rune(-1)

type Scanner struct {
	input        []byte
	offset       int
	line, column int
	scanner      func(*Scanner) (token.Token, error)
	currentChar  rune
	triviaTypes  []token.Type
	lastErr      error
}

// Init initializes the scanner.
//
// Call Init before scanning tokens with NextToken.
// Scanner middleware must be registered via UseScanner BEFORE Init.
func (sc *Scanner) Init(input []byte, opts ...func(*config)) {
	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	sc.triviaTypes = cfg.withTriviaTypes
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
	sc.currentChar = EOF
	sc.line = 0
	sc.column = -1
	sc.lastErr = nil
	sc.AdvanceChar()
}

// Err returns the last error encountered during scanning, if any.
func (sc *Scanner) Err() error {
	return sc.lastErr
}

func (sc *Scanner) CurrentChar() rune {
	return sc.currentChar
}

func (sc *Scanner) PeekChar() rune {
	if sc.offset < len(sc.input) {
		r, _ := utf8.DecodeRune(sc.input[sc.offset:])
		return r
	}
	return EOF
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
		if sc.currentChar != '\r' {
			sc.line++
			sc.column = -1
		}
	case utf8.RuneError:
		if size > 0 {
			// just an illegal character; keep going
			sc.column++
		} else {
			// reached the end of the file
			r = EOF
		}
	default:
		sc.column++
	}
	sc.currentChar = r
}

func (sc *Scanner) NextToken() token.Token {
	next := func() token.Token {
		sc.skipWhitespaces()
		line, column := sc.line, sc.column
		tok, err := sc.scanner(sc)
		if err != nil {
			tok.Type = token.ILLEGAL
			sc.lastErr = err
		}
		tok.Line = line
		tok.Column = max(0, column)
		return tok
	}
	var trivia []token.Token
	afterNewline := false
	tok := next()
triviaLoop:
	for {
		switch {
		case tok.Type == token.NEWLINE:
			afterNewline = true
		case slices.Contains(sc.triviaTypes, tok.Type):
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
	for sc.currentChar == ' ' || sc.currentChar == '\t' {
		sc.AdvanceChar()
	}
}
