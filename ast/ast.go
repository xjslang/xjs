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

type BlockStatement struct {
	Statements []Statement
}

func (node *BlockStatement) PrintTo(p *printer.Printer) {
	for _, stmt := range node.Statements {
		p.PrintIndent()
		stmt.PrintTo(p)
		p.PrintRune('\n')
	}
}

type LetStatement struct {
	Name  token.Token
	Value Expression
}

func (node *LetStatement) PrintTo(p *printer.Printer) {
	p.PrintString("let ")
	p.PrintString(node.Name.Literal)
	p.PrintString(" = ")
	node.Value.PrintTo(p)
	p.PrintRune(';')
}

type FunctionDeclaration struct {
	Name token.Token
	Body *BlockStatement
}

func (node *FunctionDeclaration) PrintTo(p *printer.Printer) {
	p.PrintString("function ")
	p.PrintString(node.Name.Literal)
	p.PrintString("() {")
	if node.Body != nil && len(node.Body.Statements) > 0 {
		p.PrintRune('\n')
		p.IncreaseIndent()
		node.Body.PrintTo(p)
		p.DecreaseIndent()
		p.PrintIndent()
	}
	p.PrintRune('}')
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
