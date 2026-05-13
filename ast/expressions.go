package ast

import "github.com/xjslang/xjs/scanner"

type InfixOperator struct {
	LeftValue  Node
	Operator   scanner.Token
	RightValue Node
}

func (node *InfixOperator) Kind() string {
	return "InfixOperator"
}

type Integer struct {
	Value string
}

func (node *Integer) Kind() string {
	return "Integer"
}

type String struct {
	Value string
}

func (node *String) Kind() string {
	return "String"
}

type Boolean struct {
	Value string
}

func (node *Boolean) Kind() string {
	return "Boolean"
}

type Ident struct {
	Value string
}

func (node *Ident) Kind() string {
	return "Ident"
}

type GroupedExpression struct {
	// keywords and delimiters
	LparenToken scanner.Token
	RparenToken scanner.Token

	Value Node
}

func (node *GroupedExpression) Kind() string {
	return "GroupedExpression"
}
