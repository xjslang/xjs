package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

// Switch tries each parser in order, restoring the parser state after each failed attempt.
func Switch[T ast.Node](p *Parser, parsers ...func(p *Parser) (T, error)) (node T, err error) {
	ss := p.scanner.(token.StatefulScanner)
	snap := ss.State()
	ct, pt := p.CurrentToken, p.PeekToken
	for _, parser := range parsers {
		if node, err = parser(p); err == nil {
			return
		}
		ss.Restore(snap)
		p.CurrentToken = ct
		p.PeekToken = pt
	}
	return
}
