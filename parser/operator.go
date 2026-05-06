package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

// RegisterInfixOperator registers an infix operator.
//
// 0: lowest precedence
// 1: +, -
// 2: *, /, %
func (p *Parser) RegisterInfixOperator(tt token.TokenType, precedence int, fn func(op token.Token, left, right ast.Expression) ast.Expression) {
	if precedence < 0 {
		panic("negative precedence")
	}
	if fn == nil {
		panic("nil function")
	}
	if p.infixOperators == nil {
		p.infixOperators = make(map[token.TokenType]infixOperator)
	}
	if _, ok := p.infixOperators[tt]; ok {
		panic("operator already registered")
	}
	p.infixOperators[tt] = infixOperator{
		precedence: precedence,
		fn:         fn,
	}
}
