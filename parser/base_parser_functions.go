package parser

import (
	"fmt"

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

func baseParseExpression(p *Parser, precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.CurrentToken.Type]
	if prefix == nil {
		p.AddError(fmt.Sprintf("no prefix parse function for %s found", p.CurrentToken.Type))
		return nil
	}
	return p.ParseRemainingExpression(prefix(), precedence)
}
