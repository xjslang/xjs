package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var ARROW = token.RegisterType("=>")

type ArrowFuncExpr struct {
	ast.BaseExpr
	Layout struct {
		Arrow token.Token
	}
	Params ast.Expr
	Body   ast.Node
}

func ParseArrowFunc(p *parser.Parser, left ast.Expr) (node *ArrowFuncExpr, err error) {
	node = &ArrowFuncExpr{Params: left}
	if node.Layout.Arrow, err = p.Expect(ARROW); err != nil {
		return
	}
	switch p.CurrentToken.Type {
	case token.LBRACE:
		if node.Body, err = js.ParseBlockStmt(p); err != nil {
			return
		}
	default:
		if node.Body, err = p.ParseExpr(); err != nil {
			return
		}
	}
	return
}

func PrintArrowFunc(pr *printer.Printer, node *ArrowFuncExpr) error {
	var printParams func(ast.Expr) error
	printParams = func(n ast.Expr) error {
		switch v := n.(type) {
		case *js.Variable:
			pr.Print(v)
		case *js.GroupExpr:
			pr.Print(v.Layout.Lparen)
			if err := printParams(v.Value); err != nil {
				return err
			}
			pr.Print(v.Layout.Rparen)
		case *SequenceExpr:
			pr.Print(v.Layout.Lparen)
			pr.IncreaseIndent()
			for i, val := range v.Values {
				if i > 0 {
					pr.Print(",")
					pr.Space()
				}
				pr.Print(val)
			}
			pr.DecreaseIndent()
			pr.Print(v.Layout.Rparen)
		default:
			// TODO: `({ a }) => {}` cannot be printed, despite being a valid expression
			return pr.Error("invalid parameter")
		}
		return nil
	}
	if err := printParams(node.Params); err != nil {
		return err
	}
	pr.Space().Print(node.Layout.Arrow)
	pr.Space().Print(node.Body)
	return nil
}
