package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func (p *Parser) UseStmtParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.stmtParser
	if next == nil {
		next = defaultStmtParser
	}
	p.stmtParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func (p *Parser) UseExprParser(parser func(p *Parser, next func() (ast.Node, error)) (ast.Node, error)) {
	next := p.exprParser
	if next == nil {
		next = defaultExprParser
	}
	p.exprParser = func(p *Parser) (ast.Node, error) {
		return parser(p, func() (ast.Node, error) {
			return next(p)
		})
	}
}

func defaultStmtParser(p *Parser) (ast.Node, error) {
	switch p.CurrentToken.Type {
	case token.LET:
		return ParseLetStmt(p)
	case token.FUNCTION:
		return ParseFuncDecl(p)
	default:
		return ParseExprStmt(p)
	}
}

func defaultExprParser(p *Parser) (val ast.Node, err error) {
	parseInfixCall := func(val ast.Node) (node *ast.CallExpr, err error) {
		node = &ast.CallExpr{
			Function: val,
		}
		if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
			return
		}
		if p.CurrentToken.Type != token.RPAREN {
			for {
				val, err := p.ParseExpr()
				if err != nil {
					return nil, err
				}
				node.Arguments = append(node.Arguments, val)
				if p.CurrentToken.Type == token.RPAREN {
					break
				}
				if _, err := p.Expect(token.COMMA); err != nil {
					return nil, err
				}
			}
		}
		if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
			return nil, err
		}
		return
	}

	if val, err = p.parseValue(); err != nil {
		return
	}
	op := p.CurrentToken
	for p.isOperator(op) {
		if op.Type == token.LPAREN {
			if val, err = parseInfixCall(val); err != nil {
				return
			}
		} else {
			binaryVal := &ast.BinaryExpr{
				LeftValue: val,
				Operator:  op,
			}
			if binaryVal.RightValue, err = ParseRemainingExpr(p); err != nil {
				return
			}
			val = binaryVal
		}
		op = p.CurrentToken
	}
	return
}
