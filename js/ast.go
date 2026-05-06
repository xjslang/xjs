package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type InfixOperator struct {
	LeftValue  ast.Expression
	Operator   token.Token
	RightValue ast.Expression
}

func (node *InfixOperator) PrintTo(p *printer.Printer) {
	node.LeftValue.PrintTo(p)
	p.PrintRune(' ')
	p.PrintString(node.Operator.Type.String())
	p.PrintRune(' ')
	node.RightValue.PrintTo(p)
}

type LetStatement struct {
	Name  token.Token
	Value ast.Expression
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

type BlockStatement struct {
	Statements []ast.Statement
}

func (node *BlockStatement) PrintTo(p *printer.Printer) {
	for _, stmt := range node.Statements {
		p.PrintIndent()
		stmt.PrintTo(p)
		p.PrintRune('\n')
	}
}
