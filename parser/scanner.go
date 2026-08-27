package parser

import (
	"errors"
	"unicode/utf8"

	"github.com/xjslang/xjs/token"
)

const EOF = rune(-1)

type ScannerState struct {
	Offset      int
	CurrentSize int
	CurrentChar rune
}

type Scanner struct {
	input       string
	offset      int
	scanner     func(*Scanner) token.Token
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

func (s *Scanner) ForkFrom(offset int) *Scanner {
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

func (s *Scanner) Fork() *Scanner {
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

func (s *Scanner) Apply(from *Scanner) {
	s.offset = from.offset
	s.currentSize = from.currentSize
	s.currentChar = from.currentChar
}

func (s *Scanner) State() ScannerState {
	return ScannerState{
		Offset:      s.offset,
		CurrentSize: s.currentSize,
		CurrentChar: s.currentChar,
	}
}

func (s *Scanner) Restore(st ScannerState) {
	s.offset = st.Offset
	s.currentSize = st.CurrentSize
	s.currentChar = st.CurrentChar
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

	tok := s.scanner(s)
	tok.LeadingTrivia = leadingTrivia
	tok.Offset = startOffset
	if triviaErr != nil {
		if e, ok := errors.AsType[Error](triviaErr); ok {
			tok.Type = e.Type
		} else {
			panic(triviaErr)
		}
	}
	return tok
}
