package js

import (
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

var blockScope = parser.RegisterScope()

func ParseProgram(p *parser.Parser) (*BlockStatement, error) {
	result := &BlockStatement{}
	for p.CurrentToken.Type != token.EOF {
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

func ParseLetStatement(p *parser.Parser) (*LetStatement, error) {
	stmt := &LetStatement{}
	p.AdvanceToken() // consume let
	ident := p.CurrentToken
	if err := p.Expect(token.IDENT); err != nil {
		return nil, err
	}
	stmt.Name = ident
	if err := p.Expect(token.ASSIGN); err != nil {
		return nil, err
	}
	val, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	stmt.Value = val
	if err := ExpectSemi(p); err != nil {
		return nil, err
	}
	return stmt, nil
}

func ParseFunctionDeclaration(p *parser.Parser) (*FunctionDeclaration, error) {
	stmt := &FunctionDeclaration{}
	p.AdvanceToken() // consume function
	ident := p.CurrentToken
	if err := p.Expect(token.IDENT); err != nil {
		return nil, err
	}
	stmt.Name = ident
	if err := p.Expect(token.LPAREN); err != nil {
		return nil, err
	}
	if err := p.Expect(token.RPAREN); err != nil {
		return nil, err
	}
	if err := p.Expect(token.LBRACE); err != nil {
		return nil, err
	}
	stmt.Body = ParseBlockStatement(p)
	if err := p.Expect(token.RBRACE); err != nil {
		return nil, err
	}
	return stmt, nil
}

func ParseBlockStatement(p *parser.Parser) *BlockStatement {
	p.EnterScope(blockScope)
	defer p.ExitScope(blockScope)
	bodyStmt := &BlockStatement{}
	for p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE {
		stmt, err := p.ParseStatement()
		if err != nil {
			AdvanceToStatementEnd(p)
			continue
		}
		bodyStmt.Statements = append(bodyStmt.Statements, stmt)
	}
	return bodyStmt
}
