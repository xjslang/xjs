package printer

import "github.com/xjslang/xjs/ast"

func PrintProgram(p *Printer, node *ast.Program) {
	for _, stmt := range node.Stmts {
		p.Print(stmt)
	}
	p.Print(node.EOFToken)
}

func PrintBlock(p *Printer, node *ast.Block) {
	p.Print(node.LbraceToken)
	p.IncreaseIndent()
	for _, stmt := range node.Stmts {
		p.Print(stmt)
	}
	// RBRACE is a special token, since the "leading trivia"
	// must be printed "before" indentation level decreases
	p.PrintTrivia(node.RbraceToken.LeadingTrivia)
	p.DecreaseIndent()
	p.LnPrint(node.RbraceToken.Literal)
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

func PrintParenExpr(p *Printer, node *ast.ParenExpr) {
	p.Print(node.LparenToken, node.Value, node.RparenToken)
}
