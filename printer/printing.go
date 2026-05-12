package printer

import "github.com/xjslang/xjs/ast"

func PrintBlock(p *Printer, node *ast.Block) {
	for _, stmt := range node.Statements {
		p.PrintIndent()
		p.Print(stmt)
		p.PrintRune('\n')
	}
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
	p.PrintString("() {")
	if node.Body != nil && len(node.Body.Statements) > 0 {
		p.PrintRune('\n')
		p.IncreaseIndent()
		p.Print(node.Body)
		p.DecreaseIndent()
		p.PrintIndent()
	}
	p.PrintRune('}')
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
