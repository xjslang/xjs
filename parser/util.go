package parser

import (
	"github.com/xjslang/xjs/ast"
)

// Switch tries each parser in order, restoring the parser state after each failed attempt.
func Switch[T ast.Node](p *Parser, parsers ...func(p *Parser) (T, error)) (node T, err error) {
	offs, curSize, curChar := p.scanner.offset, p.scanner.currentSize, p.scanner.currentChar
	curTok, peekTok := p.CurrentToken, p.PeekToken
	for _, parser := range parsers {
		if node, err = parser(p); err == nil {
			return
		}
		p.scanner.offset = offs
		p.scanner.currentSize = curSize
		p.scanner.currentChar = curChar
		p.CurrentToken = curTok
		p.PeekToken = peekTok
	}
	return
}

func IsLetter(r rune) bool {
	return r >= 'a' && r <= 'z' ||
		r >= 'A' && r <= 'Z' ||
		r == '_' || r == '$' ||
		r == '\u200c' || // ZWNJ
		r == '\u200d' // ZWJ
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isHexDigit(r rune) bool {
	return isDigit(r) || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}
