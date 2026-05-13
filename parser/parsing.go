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

func ParseGroupedExpression(p *Parser) (*ast.GroupedExpression, error) {
	node := &ast.GroupedExpression{}
	node.LparenToken = p.CurrentToken
	if err := p.Expect(scanner.LPAREN); err != nil {
		return nil, err
	}
	val, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	node.Value = val
	node.RparenToken = p.CurrentToken
	if err := p.Expect(scanner.RPAREN); err != nil {
		return nil, err
	}
	return node, nil
}

func ParseLet(p *Parser) (*ast.Let, error) {
	node := &ast.Let{}
	node.LetToken = p.CurrentToken
	if err := p.Expect(scanner.LET); err != nil {
		return nil, err
	}
	node.Name = p.CurrentToken
	if err := p.Expect(scanner.IDENT); err != nil {
		return nil, err
	}
	node.AssignToken = p.CurrentToken
	if err := p.Expect(scanner.ASSIGN); err != nil {
		return nil, err
	}
	val, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	node.Value = val
	node.SemiToken = p.CurrentToken
	if err := ExpectSemi(p); err != nil {
		return nil, err
	}
	return node, nil
}

func ParseFunction(p *Parser) (*ast.Function, error) {
	node := &ast.Function{}
	node.FunctionToken = p.CurrentToken
	if err := p.Expect(scanner.FUNCTION); err != nil {
		return nil, err
	}
	node.Name = p.CurrentToken
	if err := p.Expect(scanner.IDENT); err != nil {
		return nil, err
	}
	node.LparenToken = p.CurrentToken
	if err := p.Expect(scanner.LPAREN); err != nil {
		return nil, err
	}
	node.RparenToken = p.CurrentToken
	if err := p.Expect(scanner.RPAREN); err != nil {
		return nil, err
	}
	body, err := ParseBlock(p)
	if err != nil {
		return nil, err
	}
	node.Body = body
	return node, nil
}

func ParseBlock(p *Parser) (*ast.Block, error) {
	p.EnterScope(blockScope)
	defer p.ExitScope(blockScope)
	node := &ast.Block{}
	node.LbraceToken = p.CurrentToken
	if err := p.Expect(scanner.LBRACE); err != nil {
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
	node.RbraceToken = p.CurrentToken
	if err := p.Expect(scanner.RBRACE); err != nil {
		return nil, err
	}
	return node, nil
}

func ExpectSemi(p *Parser) error {
	if advanceSemi(p) {
		return nil
	}
	msg := "Expected statement terminator"
	p.AddError(msg)
	return errors.New(msg)
}

func AdvanceToStatementEnd(p *Parser) {
	for !advanceSemi(p) {
		p.AdvanceToken()
	}
}

func advanceSemi(p *Parser) bool {
	if p.CurrentToken.Type == scanner.SEMICOLON {
		p.AdvanceToken()
		return true
	}
	if p.CurrentToken.Type == scanner.EOF || p.CurrentToken.AfterNewline {
		return true
	}
	if p.InScope(blockScope) && p.CurrentToken.Type == scanner.RBRACE {
		return true
	}
	return false
}
