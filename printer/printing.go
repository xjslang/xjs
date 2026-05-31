package printer

import "github.com/xjslang/xjs/ast"

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
