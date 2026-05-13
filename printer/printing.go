package printer

import "github.com/xjslang/xjs/ast"

func PrintProgram(p *Printer, node *ast.Program) {
	for _, stmt := range node.Statements {
		p.Print(stmt)
		p.PrintRune('\n')
	}
}

func PrintBlock(p *Printer, node *ast.Block) {
	p.PrintString("{\n")
	p.IncreaseIndent()
	for _, stmt := range node.Statements {
		p.PrintIndent()
		p.Print(stmt)
		p.PrintRune('\n')
	}
	p.DecreaseIndent()
	p.PrintIndent()
	p.PrintRune('}')
}

func PrintLet(p *Printer, node *ast.Let) {
	p.PrintString("let ")
	p.PrintString(node.Name.Literal)
	p.PrintString(" = ")
	p.Print(node.Value)
	p.PrintRune(';')
}

func PrintFunction(p *Printer, node *ast.Function) {
	p.PrintString("function ")
	p.PrintString(node.Name.Literal)
	p.PrintString("() ")
	p.Print(node.Body)
}

func PrintInfixOperator(p *Printer, node *ast.InfixOperator) {
	p.Print(node.LeftValue)
	p.PrintRune(' ')
	p.PrintString(node.Operator.Type.String())
	p.PrintRune(' ')
	p.Print(node.RightValue)
}

func PrintGroupedExpression(p *Printer, node *ast.GroupedExpression) {
	p.PrintRune('(')
	p.Print(node.Value)
	p.PrintRune(')')
}
