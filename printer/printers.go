package printer

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
)

func IdentifierPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*ast.Identifier); ok {
		p.PrintString(node.Value)
		return
	}
	next()
}

func InfixOperatorPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*js.InfixOperator); ok {
		p.Print(node.LeftValue)
		p.PrintRune(' ')
		p.PrintString(node.Operator.Type.String())
		p.PrintRune(' ')
		p.Print(node.RightValue)
		return
	}
	next()
}

func IntegerLiteralPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*ast.IntegerLiteral); ok {
		p.PrintString(node.Value)
		return
	}
	next()
}

func StringLiteralPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*ast.StringLiteral); ok {
		p.PrintString(node.Value)
		return
	}
	next()
}

func BooleanLiteralPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*ast.BooleanLiteral); ok {
		p.PrintString(node.Value)
		return
	}
	next()
}

func GroupedExpressionPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*ast.GroupedExpression); ok {
		p.PrintRune('(')
		p.Print(node.Value)
		p.PrintRune(')')
		return
	}
	next()
}

func LetPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*js.LetStatement); ok {
		p.PrintString("let ")
		p.PrintString(node.Name.Literal)
		p.PrintString(" = ")
		p.Print(node.Value)
		p.PrintRune(';')
		return
	}
	next()
}

func BlockPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*js.BlockStatement); ok {
		for _, stmt := range node.Statements {
			p.PrintIndent()
			p.Print(stmt)
			p.PrintRune('\n')
		}
		return
	}
	next()
}

func FunctionPrinter(p *Printer, node ast.Node, next func()) {
	if node, ok := node.(*js.FunctionDeclaration); ok {
		p.PrintString("function ")
		p.PrintString(node.Name.Literal)
		p.PrintString("() {")
		if node.Body != nil && len(node.Body.Statements) > 0 {
			p.PrintRune('\n')
			p.IncreaseIndent()
			p.Print(node.Body)
			p.DecreaseIndent()
			p.PrintIndent()
		}
		p.PrintRune('}')
		return
	}
	next()
}
