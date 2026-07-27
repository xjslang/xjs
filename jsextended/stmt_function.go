package jsextended

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type Param struct {
	ast.BaseNode
	Layout struct {
		Assign token.Token
	}
	Pattern ast.Node
	Default ast.Expr
}

type RestParam struct {
	ast.BaseNode
	Layout struct {
		Spread token.Token
	}
	Name *js.Ident
}

type FunctionParams struct {
	ast.BaseNode
	Layout struct {
		Lparen token.Token
		Rparen token.Token
	}
	Params []ast.Node
}

type FunctionDecl struct {
	ast.BaseDecl
	Layout struct {
		Function token.Token
	}
	Name   *js.Ident
	Params *FunctionParams
	Body   *js.BlockStmt
}

func ParseFunctionDecl(p *parser.Parser) (node *FunctionDecl, err error) {
	node = &FunctionDecl{}
	if node.Layout.Function, err = p.Expect(js.FUNCTION); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	if node.Params, err = ParseFunctionParams(p); err != nil {
		return
	}
	if node.Body, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	return node, nil
}

func ParseFunctionParams(p *parser.Parser) (node *FunctionParams, err error) {
	node = &FunctionParams{}
	if node.Layout.Lparen, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	for p.CurrentToken.Type != token.RPAREN {
		var param ast.Node
		if p.CurrentToken.Type == SPREAD {
			if param, err = parseRestParam(p); err != nil {
				return
			}
			node.Params = append(node.Params, param)
			break
		} else if param, err = parseParam(p); err != nil {
			return
		}
		node.Params = append(node.Params, param)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
		p.AdvanceToken()
	}
	if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	return
}

func parseParam(p *parser.Parser) (param *Param, err error) {
	param = &Param{}
	switch p.CurrentToken.Type {
	case token.LBRACE:
		if param.Pattern, err = ParseObjExpr(p); err != nil {
			return
		}
	case token.LBRACKET:
		if param.Pattern, err = ParseArrayExpr(p); err != nil {
			return
		}
	default:
		if param.Pattern, err = js.ParseIdent(p); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == token.ASSIGN {
		param.Layout.Assign = p.CurrentToken
		p.AdvanceToken()
		if param.Default, err = p.ParseExpr(); err != nil {
			return
		}
	}
	return
}

func parseRestParam(p *parser.Parser) (param *RestParam, err error) {
	param = &RestParam{}
	param.Layout.Spread = p.CurrentToken
	p.AdvanceToken()
	if param.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	return
}

func PrintFunctionDecl(pr *printer.Printer, node *FunctionDecl) (err error) {
	pr.Line().Print(node.Layout.Function)
	pr.Space().Print(node.Name)
	if err = PrintFunctionParams(pr, node.Params); err != nil {
		return err
	}
	pr.Space().Print(node.Body)
	return
}

func PrintFunctionParams(pr *printer.Printer, node *FunctionParams) error {
	pr.Print(node.Layout.Lparen)
	pr.IncreaseIndent()
	for i, param := range node.Params {
		if i > 0 {
			pr.Print(",")
			pr.Space()
		}
		switch v := param.(type) {
		case *Param:
			pr.Print(v.Pattern)
			if v.Default != nil {
				pr.Space().Print(v.Layout.Assign)
				pr.Space().Print(v.Default)
			}
		case *RestParam:
			pr.Print(v.Layout.Spread, v.Name)
		default:
			return pr.Error("param expected")
		}
	}
	pr.DecreaseIndent()
	pr.Print(node.Layout.Rparen)
	return nil
}
