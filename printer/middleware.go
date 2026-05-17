package printer

import "github.com/xjslang/xjs/ast"

func (p *Printer) UsePrinter(printer func(c *Printer, node ast.Node, next func())) {
	print := p.printer
	if p.printer == nil {
		print = p.defaultPrinter
	}
	p.printer = func(node ast.Node) {
		printer(p, node, func() {
			print(node)
		})
	}
}

func (p *Printer) defaultPrinter(node ast.Node) {
	switch node := node.(type) {
	case *ast.Program:
		PrintProgram(p, node)
	case *ast.Block:
		PrintBlock(p, node)
	case *ast.Let:
		PrintLet(p, node)
	case *ast.Function:
		PrintFunction(p, node)
	case *ast.GroupedExpression:
		PrintGroupedExpression(p, node)
	case *ast.InfixOperator:
		PrintInfixOperator(p, node)
	case *ast.Ident:
		p.PrintToken(node.Value)
	case *ast.Integer:
		p.PrintToken(node.Value)
	case *ast.String:
		p.PrintToken(node.Value)
	case *ast.Boolean:
		p.PrintToken(node.Value)
	default:
		p.PrintString("<" + node.Kind() + ">")
	}
}
