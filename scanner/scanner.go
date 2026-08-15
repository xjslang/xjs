package scanner

import (
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

const EOF = rune(-1)

type Scanner struct {
	input       string
	offset      int
	scanner     func(*Scanner) (token.Token, error)
	currentSize int
	currentChar rune
}

func (s *Scanner) init(input string) {
	s.input = input
	if s.scanner == nil {
		s.scanner = defaultScanner
	}
	s.Reset()
}

func (s *Scanner) ForkFrom(offset int) token.Scanner {
	newS := &Scanner{
		input:       s.input,
		offset:      offset,
		scanner:     s.scanner,
		currentSize: 0,
		currentChar: EOF,
	}
	if newS.scanner == nil {
		newS.scanner = defaultScanner
	}
	newS.AdvanceChar()
	return newS
}

func (s *Scanner) Fork() token.Scanner {
	newS := &Scanner{
		input:       s.input,
		offset:      s.offset,
		currentSize: s.currentSize,
		currentChar: s.currentChar,
	}
	newS.scanner = s.scanner
	if newS.scanner == nil {
		newS.scanner = defaultScanner
	}
	return newS
}

func (s *Scanner) Apply(from token.Scanner) {
	switch v := from.(type) {
	case *Scanner:
		s.offset = v.offset
		s.currentSize = v.currentSize
		s.currentChar = v.currentChar
	default:
		panic("*Scanner expected")
	}
}

func (s *Scanner) Reset() {
	if s.scanner == nil {
		s.scanner = defaultScanner
	}
	s.offset = 0
	s.currentSize = 0
	s.currentChar = EOF
	s.AdvanceChar()
}

func (s *Scanner) CurrentChar() rune {
	return s.currentChar
}

func (s *Scanner) PeekChar() rune {
	if s.offset < len(s.input) {
		r, _ := utf8.DecodeRuneInString(s.input[s.offset:])
		return r
	}
	return EOF
}

func (s *Scanner) AdvanceChar() {
	r, size := utf8.DecodeRuneInString(s.input[s.offset:])
	s.offset += size
	if size == 0 && r == utf8.RuneError {
		r = EOF
	}
	s.currentChar, s.currentSize = r, size
}

func (s *Scanner) NextToken() token.Token {
	ofs0 := s.offset - s.currentSize
	triviaErr := scanTrivia(s)
	leadingTrivia := string(s.input[ofs0 : s.offset-s.currentSize])
	startOffset := s.offset - s.currentSize

	tok, err := s.scanner(s)
	tok.LeadingTrivia = leadingTrivia
	tok.Offset = startOffset
	if triviaErr != nil || err != nil {
		tok.Type = token.ILLEGAL
	}
	return tok
}
