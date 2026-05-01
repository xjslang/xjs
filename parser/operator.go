package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

// RegisterInfixOperator registers an infix operator.
// It returns an error if the operator is already registered.
//
// 0: lowest precedence
// 1: +, -
// 2: *, /, %
func (p *Parser) RegisterInfixOperator(tt token.TokenType, precedence int, fn func(op token.Token, left, right ast.Expression) ast.Expression) error {
	if precedence < 0 {
		return errors.New("negative precedence")
	}
	if fn == nil {
		return errors.New("nil function")
	}
	if p.infixOperators == nil {
		p.infixOperators = make(map[token.TokenType]infixOperator)
	}
	if _, ok := p.infixOperators[tt]; ok {
		return errors.New("operator already registered")
	}
	p.infixOperators[tt] = infixOperator{
		precedence: precedence,
		fn:         fn,
	}
	return nil
}

func defaultInfixOperator(op token.Token, left, right ast.Expression) ast.Expression {
	return &ast.InfixOperator{LeftValue: left, Operator: op, RightValue: right}
}
