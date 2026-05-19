package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) RegisterInfixOperator(typ token.Type, precedence int, fn func(op token.Token, left, right ast.Node) ast.Node) {
	if precedence < 0 {
		panic("negative precedence")
	}
	if fn == nil {
		panic("nil function")
	}
	if p.infixOperators == nil {
		p.infixOperators = make(map[token.Type]infixOperator)
	}
	if _, ok := p.infixOperators[typ]; ok {
		panic("operator already registered")
	}
	p.infixOperators[typ] = infixOperator{
		precedence: precedence,
		fn:         fn,
	}
}
