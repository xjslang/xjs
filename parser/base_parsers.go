package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

// Parses any kind of statement
func baseParseStatement(p *Parser) ast.Statement {
	switch p.CurrentToken.Type {
	case token.LET:
		return p.parseLetStatement(p)
	case token.FUNCTION:
		return p.parseFunctionDeclaration()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Parses variable declarations
func baseParseLetStatement(p *Parser) *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.CurrentToken}

	if !p.ExpectToken(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}

	if p.PeekToken.Type == token.ASSIGN {
		p.NextToken() // consume =
		p.NextToken() // move to value
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.PeekToken.Type == token.SEMICOLON {
		p.NextToken()
	}

	return stmt
}
