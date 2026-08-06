package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type InfixTemplateOp struct {
	ast.BaseExpr
	Tag      ast.Expr
	Template token.Token
}

func ParseInfixTemplateOp(p *parser.Parser, left ast.Expr) (node *InfixTemplateOp, err error) {
	node = &InfixTemplateOp{Tag: left}
	if node.Template, err = p.Expect(token.STRING); err != nil {
		return
	}
	if len(node.Template.Literal) == 0 || node.Template.Literal[0] != '`' {
		err = p.ErrorAt(node.Template, "template literal expected")
		return
	}
	return
}

func PrintInfixTemplateOp(pr *printer.Printer, node *InfixTemplateOp) error {
	pr.Print(node.Tag, node.Template)
	return nil
}
