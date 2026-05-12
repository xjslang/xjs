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

type Integer struct {
	Value string
}

func (node *Integer) Type() string {
	return "Integer"
}

type String struct {
	Value string
}

func (node *String) Type() string {
	return "String"
}

type Boolean struct {
	Value string
}

func (node *Boolean) Type() string {
	return "Boolean"
}

type Ident struct {
	Value string
}

func (node *Ident) Type() string {
	return "Ident"
}

type GroupedExpression struct {
	Value Node
}

func (node *GroupedExpression) Type() string {
	return "GroupedExpression"
}
