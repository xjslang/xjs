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
