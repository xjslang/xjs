package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
)

var blockScope = RegisterScope()

func ParseProgram(p *Parser) (*ast.Program, error) {
	result := &ast.Program{}
	for p.CurrentToken.Type != scanner.EOF {
		stmt, err := p.ParseStatement()
		if err != nil {
			AdvanceToStatementEnd(p)
			continue
		}
		result.Statements = append(result.Statements, stmt)
	}
	if errors := p.Errors(); len(errors) > 0 {
		return result, errors
	}
	return result, nil
}

func ParseGroupedExpression(p *Parser) (node *ast.GroupedExpression, err error) {
	node = &ast.GroupedExpression{}
	if node.LparenToken, err = p.Expect(scanner.LPAREN); err != nil {
		return nil, err
	}
	node.Value, err = p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if node.RparenToken, err = p.Expect(scanner.RPAREN); err != nil {
		return nil, err
	}
	return node, nil
}

func ParseLet(p *Parser) (node *ast.Let, err error) {
	node = &ast.Let{}
	if node.LetToken, err = p.Expect(scanner.LET); err != nil {
		return nil, err
	}
	if node.Name, err = p.Expect(scanner.IDENT); err != nil {
		return nil, err
	}
	if node.AssignToken, err = p.Expect(scanner.ASSIGN); err != nil {
		return nil, err
	}
	node.Value, err = p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if node.SemiToken, err = ExpectSemi(p); err != nil {
		return nil, err
	}
	return node, nil
}

func ParseFunction(p *Parser) (node *ast.Function, err error) {
	node = &ast.Function{}
	if node.FunctionToken, err = p.Expect(scanner.FUNCTION); err != nil {
		return nil, err
	}
	if node.Name, err = p.Expect(scanner.IDENT); err != nil {
		return nil, err
	}
	if node.LparenToken, err = p.Expect(scanner.LPAREN); err != nil {
		return nil, err
	}
	if node.RparenToken, err = p.Expect(scanner.RPAREN); err != nil {
		return nil, err
	}
	body, err := ParseBlock(p)
	if err != nil {
		return nil, err
	}
	node.Body = body
	return node, nil
}

func ParseBlock(p *Parser) (node *ast.Block, err error) {
	p.EnterScope(blockScope)
	defer p.ExitScope(blockScope)
	node = &ast.Block{}
	if node.LbraceToken, err = p.Expect(scanner.LBRACE); err != nil {
		return nil, err
	}
	for p.CurrentToken.Type != scanner.EOF && p.CurrentToken.Type != scanner.RBRACE {
		stmt, err := p.ParseStatement()
		if err != nil {
			AdvanceToStatementEnd(p)
			continue
		}
		node.Statements = append(node.Statements, stmt)
	}
	if node.RbraceToken, err = p.Expect(scanner.RBRACE); err != nil {
		return nil, err
	}
	return node, nil
}

func ExpectSemi(p *Parser) (scanner.Token, error) {
	tok := p.CurrentToken
	if tok.Type == scanner.SEMICOLON {
		p.AdvanceToken()
		return tok, nil
	}
	if tok.Type == scanner.EOF || tok.AfterNewline {
		tok = scanner.Token{Type: scanner.SEMICOLON, Literal: scanner.SEMICOLON.String(), Position: tok.Position}
		return tok, nil
	}
	if p.InScope(blockScope) && tok.Type == scanner.RBRACE {
		tok = scanner.Token{Type: scanner.SEMICOLON, Literal: scanner.SEMICOLON.String(), Position: tok.Position}
		return tok, nil
	}
	msg := "Expected statement terminator"
	p.AddError(msg)
	return tok, errors.New(msg)
}

func AdvanceToStatementEnd(p *Parser) {
	for {
		typ := p.CurrentToken.Type
		if typ == scanner.SEMICOLON {
			p.AdvanceToken()
			break
		}
		if typ == scanner.EOF || p.CurrentToken.AfterNewline ||
			p.InScope(blockScope) && typ == scanner.RBRACE {
			break
		}
		p.AdvanceToken()
	}
}
