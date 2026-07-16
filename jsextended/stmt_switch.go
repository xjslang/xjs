package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	SWITCH  = token.RegisterType("switch")
	CASE    = token.RegisterType("case")
	DEFAULT = token.RegisterType("default")
)

type SwitchStmt struct {
	ast.BaseStmt
	Layout struct {
		Switch token.Token
		Lparen token.Token
		Rparen token.Token
		Lbrace token.Token
		Rbrace token.Token
	}
	Expr    ast.Expr
	Clauses []ast.Stmt
}

type SwitchCaseStmt struct {
	ast.BaseStmt
	Layout struct {
		Case  token.Token
		Colon token.Token
	}
	Expr  ast.Expr
	Stmts []ast.Stmt
}

type SwitchDefaultStmt struct {
	ast.BaseStmt
	Layout struct {
		Default token.Token
		Colon   token.Token
	}
	Stmts []ast.Stmt
}

func ParseSwitchStmt(p *parser.Parser) (node *SwitchStmt, err error) {
	node = &SwitchStmt{}
	if node.Layout.Switch, err = p.Expect(SWITCH); err != nil {
		return
	}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	defClauses := 0
clausesLoop:
	for {
		var stmt ast.Stmt
		switch p.CurrentToken.Type {
		case CASE:
			if stmt, err = parseCaseStmt(p); err != nil {
				return
			}
		case DEFAULT:
			defClauses++
			if stmt, err = parseDefaultStmt(p); err != nil {
				return
			}
		default:
			break clausesLoop
		}
		node.Clauses = append(node.Clauses, stmt)
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	if defClauses > 1 {
		err = p.Error("multiple default clauses")
		return
	}
	return
}

func parseCaseStmt(p *parser.Parser) (node *SwitchCaseStmt, err error) {
	node = &SwitchCaseStmt{}
	node.Layout.Case = p.CurrentToken
	p.AdvanceToken() // consume CASE
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.Layout.Colon, err = p.Expect(token.COLON); err != nil {
		return
	}
	for {
		if t := p.CurrentToken.Type; t == CASE || t == DEFAULT || t == token.RBRACE || t == token.EOF {
			break
		}
		var stmt ast.Stmt
		if stmt, err = p.ParseStmt(); err != nil {
			return
		}
		node.Stmts = append(node.Stmts, stmt)
	}
	return
}

func parseDefaultStmt(p *parser.Parser) (node *SwitchDefaultStmt, err error) {
	node = &SwitchDefaultStmt{}
	node.Layout.Default = p.CurrentToken
	p.AdvanceToken() // consume DEFAULT
	if node.Layout.Colon, err = p.Expect(token.COLON); err != nil {
		return
	}
	for {
		if t := p.CurrentToken.Type; t == CASE || t == DEFAULT || t == token.RBRACE || t == token.EOF {
			break
		}
		var stmt ast.Stmt
		if stmt, err = p.ParseStmt(); err != nil {
			return
		}
		node.Stmts = append(node.Stmts, stmt)
	}
	return
}

func PrintSwitchStmt(pr *printer.Printer, node *SwitchStmt) error {
	pr.Line().Print(node.Layout.Switch)
	pr.Space().Print(node.Layout.Lparen)
	pr.Print(node.Expr, node.Layout.Rparen)
	pr.Space().Print(node.Layout.Lbrace)
	pr.IncreaseIndent()
	for _, stmt := range node.Clauses {
		switch v := stmt.(type) {
		case *SwitchCaseStmt:
			printSwitchCaseStmt(pr, v)
		case *SwitchDefaultStmt:
			printSwitchDefaultStmt(pr, v)
		}
	}
	pr.DecreaseIndent()
	pr.Line().Print(node.Layout.Rbrace)
	return nil
}

func printSwitchCaseStmt(pr *printer.Printer, node *SwitchCaseStmt) {
	pr.Line().Print(node.Layout.Case)
	pr.Space().Print(node.Expr)
	pr.Print(node.Layout.Colon)
	pr.IncreaseIndent()
	for _, stmt := range node.Stmts {
		pr.Print(stmt)
	}
	pr.DecreaseIndent()
}

func printSwitchDefaultStmt(pr *printer.Printer, node *SwitchDefaultStmt) {
	pr.Line().Print(node.Layout.Default)
	pr.Print(node.Layout.Colon)
	pr.IncreaseIndent()
	for _, stmt := range node.Stmts {
		pr.Print(stmt)
	}
	pr.DecreaseIndent()
}
