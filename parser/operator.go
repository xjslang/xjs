package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
)

func (p *Parser) RegisterInfixOperator(tt scanner.Kind, precedence int, fn func(op scanner.Token, left, right ast.Node) ast.Node) {
	if precedence < 0 {
		panic("negative precedence")
	}
	if fn == nil {
		panic("nil function")
	}
	if p.infixOperators == nil {
		p.infixOperators = make(map[scanner.Kind]infixOperator)
	}
	if _, ok := p.infixOperators[tt]; ok {
		panic("operator already registered")
	}
	p.infixOperators[tt] = infixOperator{
		precedence: precedence,
		fn:         fn,
	}
}
