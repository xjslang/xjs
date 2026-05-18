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

func PrintLet(p *Printer, node *ast.Let) {
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

func PrintFunction(p *Printer, node *ast.Function) {
	p.EnsureLine()
	p.PrintToken(node.FunctionToken)
	p.EnsureSpace()
	p.PrintToken(node.Name)
	p.PrintToken(node.LparenToken)
	p.PrintToken(node.RparenToken)
	p.EnsureSpace()
	p.PrintNode(node.Body)
}

func PrintCallExpr(p *Printer, node *ast.Call) {
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

func PrintInfixOperator(p *Printer, node *ast.InfixOperator) {
	p.PrintNode(node.LeftValue)
	p.EnsureSpace()
	p.PrintToken(node.Operator)
	p.EnsureSpace()
	p.PrintNode(node.RightValue)
}

func PrintGroupedExpression(p *Printer, node *ast.GroupedExpression) {
	p.PrintToken(node.LparenToken)
	p.PrintNode(node.Value)
	p.PrintToken(node.RparenToken)
}
