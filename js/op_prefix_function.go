package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FUNCTION = token.RegisterType("FUNCTION", "function")

type PrefixFunctionOp struct {
	ast.BaseExpr
	Layout struct {
		Function token.Token
		Multiply token.Token
	}
	IsGenerator bool
	Name        *Ident
	Params      *FunctionParams
	Body        *BlockStmt
}

func ParsePrefixFunctionOp(p *parser.Parser) (node *PrefixFunctionOp, err error) {
	node = &PrefixFunctionOp{}
	if node.Layout.Function, err = p.Expect(FUNCTION); err != nil {
		return
	}
	if node.IsGenerator = p.CurrentToken.Type == token.MULTIPLY; node.IsGenerator {
		node.Layout.Multiply = p.CurrentToken
		p.AdvanceToken()
	}
	if p.CurrentToken.Type == token.IDENT {
		node.Name = &Ident{Token: p.CurrentToken}
		p.AdvanceToken()
	}
	if node.Params, err = ParseFunctionParams(p); err != nil {
		return
	}
	if node.Body, err = ParseBlockStmt(p); err != nil {
		return
	}
	return node, nil
}

func PrintPrefixFunctionOp(pr *printer.Printer, node *PrefixFunctionOp) (err error) {
	pr.Print(node.Layout.Function)
	if node.IsGenerator {
		pr.Print(node.Layout.Multiply)
	}
	pr.Space()
	if node.Name != nil {
		pr.Print(node.Name)
	}
	if err = PrintFunctionParams(pr, node.Params); err != nil {
		return err
	}
	pr.Space().Print(node.Body)
	return
}
