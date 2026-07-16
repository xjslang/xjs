package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	TRY     = token.RegisterType("try")
	CATCH   = token.RegisterType("catch")
	FINALLY = token.RegisterType("finally")
)

type TryStmt struct {
	ast.BaseStmt
	Layout struct {
		Try     token.Token
		Catch   token.Token
		Lparen  token.Token
		Rparen  token.Token
		Finally token.Token
	}
	Try        *js.BlockStmt
	Catch      *js.BlockStmt
	CatchParam *js.Ident
	Finally    *js.BlockStmt
}

func ParseTryStmt(p *parser.Parser) (node *TryStmt, err error) {
	node = &TryStmt{}
	if node.Layout.Try, err = p.Expect(TRY); err != nil {
		return
	}
	if node.Try, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	if p.CurrentToken.Type == CATCH {
		node.Layout.Catch = p.CurrentToken
		p.AdvanceToken()
		if p.CurrentToken.Type == token.LPAREN {
			node.Layout.Lparen = p.CurrentToken
			p.AdvanceToken()
			if node.CatchParam, err = js.ParseIdent(p); err != nil {
				return
			}
			if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
				return
			}
		}
		if node.Catch, err = js.ParseBlockStmt(p); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == FINALLY {
		node.Layout.Finally = p.CurrentToken
		p.AdvanceToken()
		if node.Finally, err = js.ParseBlockStmt(p); err != nil {
			return
		}
	}
	if node.Catch == nil && node.Finally == nil {
		err = p.Error("missing catch or finally after try")
	}
	return
}

func PrintTryStmt(pr *printer.Printer, node *TryStmt) error {
	pr.Line().Print(node.Layout.Try)
	pr.Space().Print(node.Try)
	if node.Catch != nil {
		pr.Space().Print(node.Layout.Catch)
		if node.CatchParam != nil {
			pr.Space().Print(node.Layout.Lparen)
			pr.IncreaseIndent()
			pr.Print(node.CatchParam)
			pr.DecreaseIndent()
			pr.Print(node.Layout.Rparen)
		}
		pr.Space().Print(node.Catch)
	}
	if node.Finally != nil {
		pr.Space().Print(node.Layout.Finally)
		pr.Space().Print(node.Finally)
	}
	return nil
}
