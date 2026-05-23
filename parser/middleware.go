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

func defaultExprParser(p *Parser) (ast.Node, error) {
	registered := func(typ token.Type) bool {
		_, ok := p.infixOperators[typ]
		return ok
	}
	precedence := func(typ token.Type) int {
		if op, ok := p.infixOperators[typ]; ok {
			return op.precedence
		}
		return -1
	}
	parseTerm := func() (ast.Node, token.Token, error) {
		// parse val
		val, err := p.parseValue()
		if err != nil {
			return nil, token.Token{}, err
		}
		// parse op
		op := p.CurrentToken
		return val, op, nil
	}
	parseInfixCall := func(val ast.Node) (node *ast.CallExpr, err error) {
		node = &ast.CallExpr{
			Function:  val,
			Arguments: nil,
		}
		if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
			return nil, err
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
		return node, nil
	}
	var parseRightExp func(ast.Node, token.Token) (ast.Node, error)
	parseRightExp = func(v0 ast.Node, op0 token.Token) (node ast.Node, err error) {
		if op0.Type == token.LPAREN {
			return parseInfixCall(v0)
		}
		for {
			p.AdvanceToken()
			v1, op1, err := parseTerm()
			if err != nil {
				return nil, err
			}
			if precedence(op0.Type) < precedence(op1.Type) {
				v1, err = parseRightExp(v1, op1)
				if err != nil {
					return nil, err
				}
				op1 = p.CurrentToken
			}
			v0 = p.infixOperators[op0.Type].fn(op0, v0, v1)
			if precedence(op0.Type) > precedence(op1.Type) {
				return v0, nil
			}
			op0 = op1
		}
	}
	v, op, err := parseTerm()
	if err != nil {
		return nil, err
	}
	for registered(op.Type) {
		v, err = parseRightExp(v, op)
		if err != nil {
			return nil, err
		}
		op = p.CurrentToken
	}
	return v, nil
}
