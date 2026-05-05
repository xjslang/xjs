package ast

import (
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Node interface {
	PrintTo(p *printer.Printer)
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

func (node *IntegerLiteral) PrintTo(p *printer.Printer) {
	p.PrintString(node.Value)
}

// Implements Expression
type StringLiteral struct {
	Value string
}

func (node *StringLiteral) PrintTo(p *printer.Printer) {
	p.PrintString(node.Value)
}

// Implements Expression
type BooleanLiteral struct {
	Value string
}

func (node *BooleanLiteral) PrintTo(p *printer.Printer) {
	p.PrintString(node.Value)
}

// Implements Expression
type Identifier struct {
	Value string
}

func (node *Identifier) PrintTo(p *printer.Printer) {
	p.PrintString(node.Value)
}

// Implements Expression
type InfixOperator struct {
	LeftValue  Expression
	Operator   token.Token
	RightValue Expression
}

func (node *InfixOperator) PrintTo(p *printer.Printer) {
	node.LeftValue.PrintTo(p)
	p.PrintRune(' ')
	p.PrintString(node.Operator.Type.String())
	p.PrintRune(' ')
	node.RightValue.PrintTo(p)
}

// Implements Expression
type GroupedExpression struct {
	Value Expression
}

func (node *GroupedExpression) PrintTo(p *printer.Printer) {
	p.PrintRune('(')
	node.Value.PrintTo(p)
	p.PrintRune(')')
}
