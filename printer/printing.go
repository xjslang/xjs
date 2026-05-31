package printer

import "github.com/xjslang/xjs/ast"

func PrintProgram(p *Printer, node *ast.Program) {
	for _, stmt := range node.Stmts {
		p.Print(stmt)
	}
	p.Print(node.EOFToken)
}

func PrintExprStmt(p *Printer, node *ast.ExprStmt) {
	p.LnPrint(node.Expr)
	p.Print(node.SemiToken)
}

func PrintUnaryExpr(p *Printer, node *ast.UnaryExpr) {
	p.Print(node.Operator, node.Value)
}

func PrintBinaryExpr(p *Printer, node *ast.BinaryExpr) {
	p.Print(node.LeftValue)
	p.SpPrint(node.Operator, node.RightValue)
}
