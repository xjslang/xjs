package ast

import (
	"github.com/xjslang/xjs/token"
)

type InfixOperator struct {
	LeftValue  Node
	Operator   token.Token
	RightValue Node
}

func (node *InfixOperator) Type() string {
	return "InfixOperator"
}

type Call struct {
	LparenToken token.Token
	RparenToken token.Token

	Function  Node
	Arguments []Node
}

func (node *Call) Type() string {
	return "Call"
}

type Integer struct {
	Value token.Token
}

func (node *Integer) Type() string {
	return "Integer"
}

type String struct {
	Value token.Token
}

func (node *String) Type() string {
	return "String"
}

type Boolean struct {
	Value token.Token
}

func (node *Boolean) Type() string {
	return "Boolean"
}

type Ident struct {
	Value token.Token
}

func (node *Ident) Type() string {
	return "Ident"
}

type GroupedExpression struct {
	// keywords and delimiters
	LparenToken token.Token
	RparenToken token.Token

	Value Node
}

func (node *GroupedExpression) Type() string {
	return "GroupedExpression"
}
