package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type TaggedTemplateExpr struct {
	ast.BaseExpr
	Tag      ast.Expr
	Template token.Token
}

func ParseTagExpr(p *parser.Parser, left ast.Expr) (node *TaggedTemplateExpr, err error) {
	node = &TaggedTemplateExpr{Tag: left}
	if node.Template, err = p.Expect(token.STRING); err != nil {
		return
	}
	if len(node.Template.Literal) == 0 || node.Template.Literal[0] != '`' {
		err = p.ErrorAt(node.Template, "template literal expected")
		return
	}
	return
}

func PrintTagExpr(pr *printer.Printer, node *TaggedTemplateExpr) error {
	pr.Print(node.Tag, node.Template)
	return nil
}
