package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	CLASS   = token.RegisterType("CLASS", "class")
	EXTENDS = token.RegisterType("EXTENDS", "extends")
)

type ClassField struct {
	ast.BaseStmt
	Layout struct {
		Static token.Token
		Semi   token.Token
	}
	isStatic bool
	Key      ast.Node
	Value    ast.Expr
	Default  ast.Expr
}
type ClassMethod struct {
	ast.BaseStmt
	Layout struct {
		Static   token.Token
		Flag     token.Token // get or set
		Multiply token.Token
	}
	isStatic    bool
	isAccessor  bool
	IsGenerator bool
	Key         ast.Node
	Params      *FunctionParams
	Body        *js.BlockStmt
}

type ClassInitializer struct {
	ast.BaseStmt
	Layout struct {
		Static token.Token
	}
	Body *js.BlockStmt
}

type ClassStmt struct {
	ast.BaseDecl
	Layout struct {
		Class   token.Token
		Extends token.Token
		Lbrace  token.Token
		Rbrace  token.Token
	}
	Name    *js.Ident
	Extends ast.Expr
	Members []ast.Stmt
}

func ParseClassStmt(p *parser.Parser) (node *ClassStmt, err error) {
	node = &ClassStmt{}
	if node.Layout.Class, err = p.Expect(CLASS); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	if p.CurrentToken.Type == EXTENDS {
		node.Layout.Extends = p.CurrentToken
		p.AdvanceToken()
		if node.Extends, err = p.ParseExpr(); err != nil {
			return
		}
	}
	if node.Layout.Lbrace, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	for {
		for p.CurrentToken.Type == token.SEMICOLON {
			p.AdvanceToken()
		}
		if p.CurrentToken.Type == token.RBRACE || p.CurrentToken.Type == token.EOF {
			break
		}
		var member ast.Stmt
		if member, err = parser.Switch(p,
			func(p *parser.Parser) (ast.Stmt, error) { return parseClassInitializer(p) },
			func(p *parser.Parser) (ast.Stmt, error) { return parseClassField(p) },
			func(p *parser.Parser) (ast.Stmt, error) { return parseClassMethod(p) },
		); err != nil {
			return
		}
		node.Members = append(node.Members, member)
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	return
}

func parseClassMethod(p *parser.Parser) (node *ClassMethod, err error) {
	node = &ClassMethod{}
	if p.CurrentToken.Literal == "static" {
		if typ := p.PeekToken.Type; typ == token.IDENT || typ == token.LBRACKET || typ == token.NUMBER || typ == token.STRING || typ == token.MULTIPLY {
			node.isStatic = true
			node.Layout.Static = p.CurrentToken
			p.AdvanceToken()
		}
	}
	if lit := p.CurrentToken.Literal; (lit == "get" || lit == "set") && p.PeekToken.Type != token.LPAREN {
		node.isAccessor = true
		node.Layout.Flag = p.CurrentToken
		p.AdvanceToken()
	}
	if node.IsGenerator = p.CurrentToken.Type == token.MULTIPLY; node.IsGenerator {
		node.Layout.Multiply = p.CurrentToken
		p.AdvanceToken()
	}
	switch p.CurrentToken.Type {
	case token.LBRACKET:
		if node.Key, err = js.ParseComputedExpr(p); err != nil {
			return
		}
	case token.STRING, token.NUMBER:
		if node.Key, err = js.ParseValue(p); err != nil {
			return
		}
	default:
		if node.Key, err = js.ParseObjKey(p); err != nil {
			return
		}
	}
	if node.Params, err = ParseFunctionParams(p); err != nil {
		return
	}
	if node.Body, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	return
}

func parseClassField(p *parser.Parser) (node *ClassField, err error) {
	node = &ClassField{}
	if p.CurrentToken.Literal == "static" {
		switch p.PeekToken.Type {
		case token.IDENT, token.NUMBER, token.STRING, token.LBRACKET, token.MULTIPLY:
			node.isStatic = true
			node.Layout.Static = p.CurrentToken
			p.AdvanceToken()
		}
	}
	switch p.CurrentToken.Type {
	case token.LBRACKET:
		if node.Key, err = js.ParseComputedExpr(p); err != nil {
			return
		}
	case token.STRING, token.NUMBER:
		if node.Key, err = js.ParseValue(p); err != nil {
			return
		}
	default:
		if node.Key, err = js.ParseObjKey(p); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.COLON {
		p.AdvanceToken()
		if node.Value, err = js.ParseRightExpr(p, token.ASSIGN.Precedence()); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.ASSIGN {
		p.AdvanceToken()
		if node.Default, err = js.ParseRightExpr(p, token.COMMA.Precedence()); err != nil {
			return
		}
	}
	if node.Layout.Semi, err = p.ExpectSemi(); err != nil {
		return
	}
	return
}

func parseClassInitializer(p *parser.Parser) (node *ClassInitializer, err error) {
	node = &ClassInitializer{}
	if node.Layout.Static, err = p.ExpectLiteral("static"); err != nil {
		return
	}
	if node.Body, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	return
}

func PrintClassStmt(pr *printer.Printer, node *ClassStmt) (err error) {
	pr.Line().Print(node.Layout.Class)
	pr.Space().Print(node.Name)
	if node.Extends != nil {
		pr.Space().Print(node.Layout.Extends)
		pr.Space().Print(node.Extends)
	}
	pr.Space().Print(node.Layout.Lbrace)
	if len(node.Members) > 0 {
		pr.IncreaseIndent()
		for _, entry := range node.Members {
			switch v := entry.(type) {
			case *ClassField:
				err = printClassField(pr, v)
			case *ClassMethod:
				err = printClassMethod(pr, v)
			case *ClassInitializer:
				err = printClassInitializer(pr, v)
			default:
				err = pr.Error("class member expected")
			}
			if err != nil {
				return
			}
		}
		pr.DecreaseIndent()
	}
	pr.Line().Print(node.Layout.Rbrace)
	return
}

func printClassField(pr *printer.Printer, node *ClassField) (err error) {
	pr.Line()
	if node.isStatic {
		pr.Print(node.Layout.Static)
	}
	pr.Space()
	if err = printMemberKey(pr, node.Key); err != nil {
		return
	}
	if node.Value != nil {
		pr.Print(":")
		pr.Space().Print(node.Value)
	}
	if node.Default != nil {
		pr.Space().Print("=")
		pr.Space().Print(node.Default)
	}
	pr.Print(node.Layout.Semi)
	return
}

func printClassMethod(pr *printer.Printer, node *ClassMethod) (err error) {
	pr.Line()
	if node.isStatic {
		pr.Print(node.Layout.Static)
		pr.Space()
	}
	if node.isAccessor {
		pr.Print(node.Layout.Flag)
	}
	if node.IsGenerator {
		pr.Space().Print(node.Layout.Multiply)
		if err = printMemberKey(pr, node.Key); err != nil {
			return
		}
	} else {
		pr.Space()
		if err = printMemberKey(pr, node.Key); err != nil {
			return
		}
	}
	if err = PrintFunctionParams(pr, node.Params); err != nil {
		return
	}
	pr.Space().Print(node.Body)
	return
}

func printClassInitializer(pr *printer.Printer, node *ClassInitializer) (err error) {
	pr.Line().Print(node.Layout.Static)
	pr.Space().Print(node.Body)
	return
}

func printMemberKey(pr *printer.Printer, node ast.Node) (err error) {
	switch v := node.(type) {
	case *js.ComputedExpr:
		pr.Print(v.Layout.Lbracket, v.Expr, v.Layout.Rbracket)
	default:
		pr.Print(v)
	}
	return
}
