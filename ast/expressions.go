package ast

import "github.com/xjslang/xjs/scanner"

type InfixOperator struct {
	LeftValue  Node
	Operator   scanner.Token
	RightValue Node
}

func (node *InfixOperator) Type() string {
	return "InfixOperator"
}

type Call struct {
	LparenToken scanner.Token
	RparenToken scanner.Token

	Function  Node
	Arguments []Node
}

func (node *Call) Type() string {
	return "Call"
}

type Integer struct {
	Value scanner.Token
}

func (node *Integer) Type() string {
	return "Integer"
}

type String struct {
	Value scanner.Token
}

func (node *String) Type() string {
	return "String"
}

type Boolean struct {
	Value scanner.Token
}

func (node *Boolean) Type() string {
	return "Boolean"
}

type Ident struct {
	Value scanner.Token
}

func (node *Ident) Type() string {
	return "Ident"
}

type GroupedExpression struct {
	// keywords and delimiters
	LparenToken scanner.Token
	RparenToken scanner.Token

	Value Node
}

func (node *GroupedExpression) Type() string {
	return "GroupedExpression"
}
