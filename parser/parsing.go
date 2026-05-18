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
			p.AdvanceToStatementEnd()
			continue
		}
		result.Statements = append(result.Statements, stmt)
	}
	result.EOFToken = p.CurrentToken
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
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
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
	var errs []error
	for p.CurrentToken.Type != scanner.EOF && p.CurrentToken.Type != scanner.RBRACE {
		stmt, err := p.ParseStatement()
		if err != nil {
			errs = append(errs, err)
			p.AdvanceToStatementEnd()
			continue
		}
		node.Statements = append(node.Statements, stmt)
	}
	if node.RbraceToken, err = p.Expect(scanner.RBRACE); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return node, nil
}
