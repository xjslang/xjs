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

func PrintLetStmt(p *Printer, node *ast.LetStmt) {
	p.LnPrint(node.LetToken)
	p.SpPrint(node.Name, node.AssignToken, node.Value)
	p.Print(node.SemiToken)
}

func PrintCallExpr(p *Printer, node *ast.CallExpr) {
	p.Print(node.Function, node.LparenToken)
	for i, arg := range node.Arguments {
		if i > 0 {
			p.Print(",")
			p.EnsureSpace()
		}
		p.Print(arg)
	}
	p.Print(node.RparenToken)
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
