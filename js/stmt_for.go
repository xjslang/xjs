package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var FOR = token.RegisterType("for")

type ForStmt struct {
	ast.BaseStmt
	Layout struct {
		For    token.Token
		Lparen token.Token
		Semi1  token.Token
		Semi2  token.Token
		Rparen token.Token
	}
	Init  ast.Stmt
	Cond  ast.Expr
	After ast.Stmt
	Then  ast.Stmt
}

func (node *ForStmt) SelfClosing() bool {
	if v, ok := node.Then.(SelfClosingStmt); ok {
		return v.SelfClosing()
	}
	return false
}

func ParseForStmt(p *parser.Parser) (node *ForStmt, err error) {
	node = &ForStmt{}
	// for
	if node.Layout.For, err = p.Expect(FOR); err != nil {
		return
	}
	// (init; cond; after)
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Init, err = parseInitClause(p); err != nil {
		return
	}
	if node.Layout.Semi1, err = p.Expect(token.SEMICOLON); err != nil {
		return
	}
	if node.Cond, err = parseCondClause(p); err != nil {
		return
	}
	if node.Layout.Semi2, err = p.Expect(token.SEMICOLON); err != nil {
		return
	}
	if node.After, err = parseAfterClause(p); err != nil {
		return
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	// then
	if node.Then, err = p.ParseStmt(); err != nil {
		return
	}
	return node, nil
}

func parseInitClause(p *parser.Parser) (node ast.Stmt, err error) {
	if p.CurrentToken.Type == token.SEMICOLON {
		// omit init clause
		return
	}
	switch p.CurrentToken.Type {
	case LET:
		n := &LetStmt{}
		n.Layout.Let = p.CurrentToken
		p.AdvanceToken()
		if n.Name, err = ParseIdent(p); err != nil {
			return
		}
		if n.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
			return
		}
		if n.Value, err = p.ParseExpr(); err != nil {
			return
		}
		node = n
	case token.IDENT:
		n := &AssignStmt{}
		if n.Name, err = ParseIdent(p); err != nil {
			return
		}
		if n.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
			return
		}
		if n.Value, err = p.ParseExpr(); err != nil {
			return
		}
		node = n
	default:
		err = p.Error("init expected")
		return
	}
	return
}

func parseCondClause(p *parser.Parser) (node ast.Expr, err error) {
	if p.CurrentToken.Type == token.SEMICOLON {
		// omit cond clause
		return
	}
	return p.ParseExpr()
}

func parseAfterClause(p *parser.Parser) (node ast.Stmt, err error) {
	if p.CurrentToken.Type == token.RPAREN {
		// omit after clause
		return
	}
	var ident *Ident
	if ident, err = ParseIdent(p); err != nil {
		return
	}
	switch p.CurrentToken.Type {
	case token.INCREMENT:
		n := &IncStmt{Name: ident}
		n.Layout.Increment = p.CurrentToken
		p.AdvanceToken()
		node = n
	case token.DECREMENT:
		n := &DecStmt{Name: ident}
		n.Layout.Decrement = p.CurrentToken
		p.AdvanceToken()
		node = n
	default:
		err = p.Error("++/-- expected")
	}
	return
}

func PrintForStmt(p *printer.Printer, node *ForStmt) {
	// for
	p.LnPrint(node.Layout.For)
	p.SpPrint(node.Layout.Lparen)
	// (init; condition; after)
	p.IncreaseIndent()
	if node.Init != nil {
		printInitClause(p, node.Init)
	}
	p.Print(node.Layout.Semi1)
	if node.Cond != nil {
		p.SpPrint(node.Cond)
	}
	p.Print(node.Layout.Semi2)
	if node.After != nil {
		printAfterClause(p, node.After)
	}
	p.DecreaseIndent()
	// then
	p.Print(node.Layout.Rparen)
	p.SpPrint(node.Then)
}

func printInitClause(p *printer.Printer, node ast.Stmt) {
	switch v := node.(type) {
	case *LetStmt:
		p.Print(v.Layout.Let)
		p.SpPrint(v.Name)
		p.SpPrint(v.Layout.Assign)
		p.SpPrint(v.Value)
	case *AssignStmt:
		p.Print(v.Name)
		p.SpPrint(v.Layout.Assign)
		p.SpPrint(v.Value)
	default:
		panic("unexpected init clause type")
	}
}

func printAfterClause(p *printer.Printer, node ast.Stmt) {
	switch v := node.(type) {
	case *IncStmt:
		p.SpPrint(v.Name)
		p.Print(v.Layout.Increment)
	case *DecStmt:
		p.SpPrint(v.Name)
		p.Print(v.Layout.Decrement)
	default:
		panic("unexpected after clause type")
	}
}
