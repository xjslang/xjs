package parser

import (
	"github.com/xjslang/xjs/ast"
)

// Switch tries each parser in order, restoring the parser state after each failed attempt.
func Switch[T ast.Node](p *Parser, parsers ...func(p *Parser) (T, error)) (node T, err error) {
	snap := p.scanner.State()
	ct, pt := p.CurrentToken, p.PeekToken
	for _, parser := range parsers {
		if node, err = parser(p); err == nil {
			return
		}
		p.scanner.Restore(snap)
		p.CurrentToken = ct
		p.PeekToken = pt
	}
	return
}
