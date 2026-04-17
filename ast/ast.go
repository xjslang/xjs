package ast

import "github.com/xjslang/xjs/printer"

type Node interface {
	PrintTo(p *printer.Printer)
}

type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Value string
}

func (id *Identifier) PrintTo(p *printer.Printer) {
	p.PrintString(id.Value)
}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) PrintTo(p *printer.Printer) {
	for _, stmt := range bs.Statements {
		p.PrintIndent()
		stmt.PrintTo(p)
		p.PrintRune('\n')
	}
}

type LetStatement struct {
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) PrintTo(p *printer.Printer) {
	p.PrintString("let ")
	ls.Name.PrintTo(p)
	p.PrintString(" = ")
	ls.Value.PrintTo(p)
	p.PrintRune(';')
}

type FunctionDeclaration struct {
	Name *Identifier
	Body *BlockStatement
}

func (fd *FunctionDeclaration) statementNode() {}

func (fd *FunctionDeclaration) PrintTo(p *printer.Printer) {
	p.PrintString("function ")
	fd.Name.PrintTo(p)
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

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) PrintTo(p *printer.Printer) {
	p.PrintString(il.Value)
}

// Implements Expression
type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) PrintTo(p *printer.Printer) {
	p.PrintString(sl.Value)
}
