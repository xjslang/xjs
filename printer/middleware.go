package printer

import "github.com/xjslang/xjs/ast"

func (p *Printer) UsePrinter(printer func(pr *Printer, node ast.Node, next func(node ast.Node))) {
	print := p.printer
	if p.printer == nil {
		print = defaultPrinter
	}
	p.printer = func(p *Printer, node ast.Node) {
		printer(p, node, func(node ast.Node) {
			print(p, node)
		})
	}
}

func defaultPrinter(p *Printer, node ast.Node) {
	switch node := node.(type) {
	case *ast.Program:
		PrintProgram(p, node)
	case *ast.Block:
		PrintBlock(p, node)
	case *ast.ExprStmt:
		PrintExprStmt(p, node)
	case *ast.LetStmt:
		PrintLetStmt(p, node)
	case *ast.FuncDecl:
		PrintFuncDecl(p, node)
	case *ast.CallExpr:
		PrintCallExpr(p, node)
	case *ast.ParenExpr:
		PrintParenExpr(p, node)
	case *ast.BinaryExpr:
		PrintBinaryExpr(p, node)
	case *ast.Ident:
		p.PrintToken(node.Value)
	case *ast.BasicLit:
		p.PrintToken(node.Value)
	default:
		p.PrintString("<" + node.Type() + ">")
	}
}
