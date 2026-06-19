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

const eof = rune(-1)

type Scanner struct {
	input        []byte
	offset       int
	line, column int
	scanner      func(*Scanner) (token.Token, error)
	currentChar  rune
	triviaTypes  []token.Type
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
	sc.currentChar = eof
	sc.line = 0
	sc.column = -1
	sc.AdvanceChar()
}

func (sc *Scanner) CurrentChar() rune {
	return sc.currentChar
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
			r = eof
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
		// TODO: (medium) Scanner.NextToken converts scanner/middleware errors into token.ILLEGAL but discards the error value entirely. With the new middleware signature returning errors, callers still have no way to observe why a token is illegal other than inspecting Literal. Consider exposing the error (e.g., NextToken returning (token.Token, error) or storing the last error on Scanner) so downstream code can surface better diagnostics.
		if err != nil {
			tok.Type = token.ILLEGAL
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
