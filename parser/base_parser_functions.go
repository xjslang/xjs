package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

func baseParseStatement(p *XJSParser) ast.Statement {
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

func baseParseExpression(p *XJSParser, precedence int) ast.Expression {
	return p.ParseRemainingExpressionWithPrecedence(
		p.ParsePrefixExpression(),
		precedence,
	)
}
