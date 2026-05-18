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

type Call struct {
	LparenToken scanner.Token
	RparenToken scanner.Token

	Function  Node
	Arguments []Node
}

func (node *Call) Kind() string {
	return "Call"
}

type Integer struct {
	Value scanner.Token
}

func (node *Integer) Kind() string {
	return "Integer"
}

type String struct {
	Value scanner.Token
}

func (node *String) Kind() string {
	return "String"
}

type Boolean struct {
	Value scanner.Token
}

func (node *Boolean) Kind() string {
	return "Boolean"
}

type Ident struct {
	Value scanner.Token
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
