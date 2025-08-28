package parser

import (
	"fmt"
	"strconv"

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

func baseParseIdentifier(p *Parser) ast.Expression {
	return &ast.Identifier{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func baseParseIntegerLiteral(p *Parser) ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.CurrentToken}
	value, err := strconv.ParseInt(p.CurrentToken.Literal, 0, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as integer", p.CurrentToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func baseParseFloatLiteral(p *Parser) ast.Expression {
	lit := &ast.FloatLiteral{Token: p.CurrentToken}
	value, err := strconv.ParseFloat(p.CurrentToken.Literal, 64)
	if err != nil {
		p.AddError(fmt.Sprintf("could not parse %q as float", p.CurrentToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

func baseParseStringLiteral(p *Parser) ast.Expression {
	return &ast.StringLiteral{Token: p.CurrentToken, Value: p.CurrentToken.Literal}
}

func baseParseBooleanLiteral(p *Parser) ast.Expression {
	return &ast.BooleanLiteral{Token: p.CurrentToken, Value: p.CurrentToken.Type == token.TRUE}
}

func baseParseNullLiteral(p *Parser) ast.Expression {
	return &ast.NullLiteral{Token: p.CurrentToken}
}

func baseParseUnaryExpression(p *Parser) ast.Expression {
	expression := &ast.UnaryExpression{
		Token:    p.CurrentToken,
		Operator: p.CurrentToken.Literal,
	}
	p.NextToken()
	expression.Right = p.ParseExpression(UNARY)
	return expression
}

func baseParseGroupedExpression(p *Parser) ast.Expression {
	p.NextToken()
	exp := p.ParseExpression(LOWEST)
	if !p.ExpectToken(token.RPAREN) {
		return nil
	}
	return exp
}

func baseParseArrayLiteral(p *Parser) ast.Expression {
	array := &ast.ArrayLiteral{Token: p.CurrentToken}
	array.Elements = p.ParseExpressionList(token.RBRACKET)
	return array
}

func baseParseObjectLiteral(p *Parser) ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.CurrentToken}
	obj.Properties = make(map[ast.Expression]ast.Expression)
	if p.PeekToken.Type == token.RBRACE {
		p.NextToken()
		return obj
	}
	p.NextToken()
	for {
		key := p.ParseExpression(LOWEST)
		if !p.ExpectToken(token.COLON) {
			return nil
		}
		p.NextToken()
		value := p.ParseExpression(LOWEST)
		obj.Properties[key] = value
		if p.PeekToken.Type != token.COMMA {
			break
		}
		p.NextToken()
		p.NextToken()
	}
	if !p.ExpectToken(token.RBRACE) {
		return nil
	}
	return obj
}

func baseParseFunctionExpression(p *Parser) ast.Expression {
	fe := &ast.FunctionExpression{Token: p.CurrentToken}
	if !p.ExpectToken(token.LPAREN) {
		return nil
	}
	fe.Parameters = p.ParseFunctionParameters()
	if !p.ExpectToken(token.LBRACE) {
		return nil
	}
	p.PushContext(FunctionContext)
	defer p.PopContext()
	fe.Body = p.ParseBlockStatement()
	return fe
}
