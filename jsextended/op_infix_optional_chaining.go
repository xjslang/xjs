package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var OPTIONAL_CHAINING = token.RegisterType("OPTIONAL_CHAINING", "?.")

type InfixOptionalChainingOp struct {
	ast.BaseExpr
	Layout struct {
		OptionalChaining token.Token
	}
	Left  ast.Expr
	Right ast.Expr
}

func ParseInfixOptionalChainingOp(p *parser.Parser, left ast.Expr) (node *InfixOptionalChainingOp, err error) {
	node = &InfixOptionalChainingOp{Left: left}
	if node.Layout.OptionalChaining, err = p.Expect(OPTIONAL_CHAINING); err != nil {
		return
	}
	// TODO: ParseInfixOptionalChainingOp delegates to js.ParseRightExpr. When the token after `?.` is `(`, ParseRightExpr parses it as a PrefixParenOp (see js.ParsePrefixParenOp), so `fn?.(a, b)` will fail at the comma (PrefixParenOp only parses a single expression and then expects `)`), and `?.[` will be treated as starting an array literal instead of optional element access. Optional chaining needs to parse the specific postfix forms (`?.ident`, `?.(args...)`, `?.[expr]`) similarly to how the core parser handles `.`, `(`, and `[` as binary postfix operators.
	if node.Right, err = js.ParseRightExpr(p, node.Layout.OptionalChaining.Type.Precedence()); err != nil {
		return
	}
	return
}

func PrintInfixOptionalChainingOp(pr *printer.Printer, node *InfixOptionalChainingOp) error {
	pr.Print(node.Left, node.Layout.OptionalChaining, node.Right)
	return nil
}
