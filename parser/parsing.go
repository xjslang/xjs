package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
)

var blockScope = RegisterScope()

func ParseProgram(p *Parser) (*ast.Block, error) {
	result := &ast.Block{}
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

func ParseLet(p *Parser) (*ast.Let, error) {
	node := &ast.Let{}
	p.AdvanceToken() // consume let
	ident := p.CurrentToken
	if err := p.Expect(scanner.IDENT); err != nil {
		return nil, err
	}
	node.Name = ident
	if err := p.Expect(scanner.ASSIGN); err != nil {
		return nil, err
	}
	val, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	node.Value = val
	if err := ExpectSemi(p); err != nil {
		return nil, err
	}
	return node, nil
}

func ParseFunction(p *Parser) (*ast.Function, error) {
	node := &ast.Function{}
	p.AdvanceToken() // consume function
	ident := p.CurrentToken
	if err := p.Expect(scanner.IDENT); err != nil {
		return nil, err
	}
	node.Name = ident
	if err := p.Expect(scanner.LPAREN); err != nil {
		return nil, err
	}
	if err := p.Expect(scanner.RPAREN); err != nil {
		return nil, err
	}
	if err := p.Expect(scanner.LBRACE); err != nil {
		return nil, err
	}
	node.Body = ParseBlock(p)
	if err := p.Expect(scanner.RBRACE); err != nil {
		return nil, err
	}
	return node, nil
}

func ParseBlock(p *Parser) *ast.Block {
	p.EnterScope(blockScope)
	defer p.ExitScope(blockScope)
	bodyStmt := &ast.Block{}
	for p.CurrentToken.Type != scanner.EOF && p.CurrentToken.Type != scanner.RBRACE {
		stmt, err := p.ParseStatement()
		if err != nil {
			AdvanceToStatementEnd(p)
			continue
		}
		bodyStmt.Statements = append(bodyStmt.Statements, stmt)
	}
	return bodyStmt
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
