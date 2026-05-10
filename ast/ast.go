package ast

type Node interface {
	Type() string
}

type Statement interface {
	Node
}
type Expression interface {
	Node
}

// Implements Expression
type IntegerLiteral struct {
	Value string
}

func (node *IntegerLiteral) Type() string {
	return "IntegerLiteral"
}

// Implements Expression
type StringLiteral struct {
	Value string
}

func (node *StringLiteral) Type() string {
	return "StringLiteral"
}

// Implements Expression
type BooleanLiteral struct {
	Value string
}

func (node *BooleanLiteral) Type() string {
	return "BooleanLiteral"
}

// Implements Expression
type Identifier struct {
	Value string
}

func (node *Identifier) Type() string {
	return "Identifier"
}

// Implements Expression
type GroupedExpression struct {
	Value Expression
}

func (node *GroupedExpression) Type() string {
	return "GroupedExpression"
}
