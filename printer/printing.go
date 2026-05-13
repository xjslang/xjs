package printer

import (
	"github.com/xjslang/xjs/ast"
)

func PrintProgram(p *Printer, node *ast.Program) {
	for _, stmt := range node.Statements {
		p.Print(stmt)
	}
}

func PrintBlock(p *Printer, node *ast.Block) {
	p.PrintTokenAfterSpace(node.LbraceToken)
	p.IncreaseIndent()
	for _, stmt := range node.Statements {
		p.Print(stmt)
	}
	// RBRACE is a special token, since the "leading trivia"
	// must be printed "before" indentation level decreases
	p.PrintTrivia(node.RbraceToken.LeadingTrivia)
	p.DecreaseIndent()
	p.EnsureLine()
	p.PrintIndent()
	p.PrintString(node.RbraceToken.Literal)
}

func PrintLet(p *Printer, node *ast.Let) {
	p.PrintTokenAfterNewline(node.LetToken)
	p.PrintTokenAfterSpace(node.Name)
	p.PrintTokenAfterSpace(node.AssignToken)
	p.Print(node.Value)
	p.PrintToken(node.SemiToken)
}

func PrintFunction(p *Printer, node *ast.Function) {
	p.PrintTokenAfterNewline(node.FunctionToken)
	p.PrintTokenAfterSpace(node.Name)
	p.PrintToken(node.LparenToken)
	p.PrintToken(node.RparenToken)
	p.Print(node.Body)
}

func PrintInfixOperator(p *Printer, node *ast.InfixOperator) {
	p.Print(node.LeftValue)
	p.PrintTokenAfterSpace(node.Operator)
	p.Print(node.RightValue)
}

func PrintGroupedExpression(p *Printer, node *ast.GroupedExpression) {
	p.PrintToken(node.LparenToken)
	p.Print(node.Value)
	p.PrintToken(node.RparenToken)
}
