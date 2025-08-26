package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func baseParseStatement(p *Parser) ast.Statement {
	switch p.CurrentToken.Type {
	case token.LET:
		return p.ParseLetStatement()
	case token.FUNCTION:
		return p.ParseFunctionStatement()
	case token.RETURN:
		return p.ParseReturnStatement()
	case token.IF:
		return p.ParseIfStatement()
	case token.WHILE:
		return p.ParseWhileStatement()
	case token.FOR:
		return p.ParseForStatement()
	case token.LBRACE:
		return p.ParseBlockStatement()
	default:
		return p.ParseExpressionStatement()
	}
}

func baseParseExpressionStatement(p *Parser) *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.CurrentToken}
	stmt.Expression = p.ParseExpression(LOWEST)

	if p.PeekToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}
