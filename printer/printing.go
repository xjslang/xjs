package printer

import "github.com/xjslang/xjs/ast"

func PrintProgram(p *Printer, node *ast.Program) {
	for _, stmt := range node.Statements {
		p.PrintNode(stmt)
	}
	p.PrintToken(node.EOFToken)
}

func PrintBlock(p *Printer, node *ast.Block) {
	p.PrintToken(node.LbraceToken)
	p.IncreaseIndent()
	for _, stmt := range node.Statements {
		p.PrintNode(stmt)
	}
	// RBRACE is a special token, since the "leading trivia"
	// must be printed "before" indentation level decreases
	p.PrintTrivia(node.RbraceToken.LeadingTrivia)
	p.DecreaseIndent()
	p.EnsureLine()
	p.PrintIndentedString(node.RbraceToken.Literal)
}

func PrintLetStmt(p *Printer, node *ast.LetStmt) {
	p.EnsureLine()
	p.PrintToken(node.LetToken)
	p.EnsureSpace()
	p.PrintToken(node.Name)
	p.EnsureSpace()
	p.PrintToken(node.AssignToken)
	p.EnsureSpace()
	p.PrintNode(node.Value)
	p.PrintToken(node.SemiToken)
}

func PrintFuncDecl(p *Printer, node *ast.FuncDecl) {
	p.EnsureLine()
	p.PrintToken(node.FunctionToken)
	p.EnsureSpace()
	p.PrintToken(node.Name)
	p.PrintToken(node.LparenToken)
	p.PrintToken(node.RparenToken)
	p.EnsureSpace()
	p.PrintNode(node.Body)
}

func PrintCallExpr(p *Printer, node *ast.CallExpr) {
	p.PrintNode(node.Function)
	p.PrintToken(node.LparenToken)
	for i, arg := range node.Arguments {
		if i > 0 {
			p.PrintIndentedString(",")
			p.EnsureSpace()
		}
		p.PrintNode(arg)
	}
	p.PrintToken(node.RparenToken)
}

func PrintBinaryExpr(p *Printer, node *ast.BinaryExpr) {
	p.PrintNode(node.LeftValue)
	p.EnsureSpace()
	p.PrintToken(node.Operator)
	p.EnsureSpace()
	p.PrintNode(node.RightValue)
}

func PrintParenExpr(p *Printer, node *ast.ParenExpr) {
	p.PrintToken(node.LparenToken)
	p.PrintNode(node.Value)
	p.PrintToken(node.RparenToken)
}
