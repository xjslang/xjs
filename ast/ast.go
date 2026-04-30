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

func (bs *BlockStatement) PrintTo(p *printer.Printer) {
	for _, stmt := range bs.Statements {
		p.PrintIndent()
		stmt.PrintTo(p)
		p.PrintRune('\n')
	}
}

type LetStatement struct {
	Name  token.Token
	Value Expression
}

func (ls *LetStatement) PrintTo(p *printer.Printer) {
	p.PrintString("let ")
	p.PrintString(ls.Name.Literal)
	p.PrintString(" = ")
	ls.Value.PrintTo(p)
	p.PrintRune(';')
}

type FunctionDeclaration struct {
	Name token.Token
	Body *BlockStatement
}

func (fd *FunctionDeclaration) PrintTo(p *printer.Printer) {
	p.PrintString("function ")
	p.PrintString(fd.Name.Literal)
	p.PrintString("() {")
	if fd.Body != nil && len(fd.Body.Statements) > 0 {
		p.PrintRune('\n')
		p.IncreaseIndent()
		fd.Body.PrintTo(p)
		p.DecreaseIndent()
		p.PrintIndent()
	}
	p.PrintRune('}')
}

// Implements Expression
type IntegerLiteral struct {
	Value string
}

func (il *IntegerLiteral) PrintTo(p *printer.Printer) {
	p.PrintString(il.Value)
}

// Implements Expression
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) PrintTo(p *printer.Printer) {
	p.PrintString(sl.Value)
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
