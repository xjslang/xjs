package parser

import (
	"errors"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

var blockScope = RegisterScope()

func ParseProgram(p *Parser) (node *ast.Program, err error) {
	node = &ast.Program{}
	for p.CurrentToken.Type != token.EOF {
		stmt, err := p.ParseStmt()
		if err != nil {
			p.AdvanceToStmtEnd()
			continue
		}
		node.Stmts = append(node.Stmts, stmt)
	}
	node.EOFToken = p.CurrentToken
	if errors := p.Errors(); len(errors) > 0 {
		err = errors
		return
	}
	return
}

func ParseExprStmt(p *Parser) (node *ast.ExprStmt, err error) {
	node = &ast.ExprStmt{}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return
}

func ParseParenExpr(p *Parser) (node *ast.ParenExpr, err error) {
	node = &ast.ParenExpr{}
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	return
}

func ParseLetStmt(p *Parser) (node *ast.LetStmt, err error) {
	node = &ast.LetStmt{}
	if node.LetToken, err = p.Expect(token.LET); err != nil {
		return
	}
	if node.Name, err = p.Expect(token.IDENT); err != nil {
		return
	}
	if node.AssignToken, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	if node.SemiToken, err = p.ExpectSemi(); err != nil {
		return
	}
	return
}

func ParseFuncDecl(p *Parser) (node *ast.FuncDecl, err error) {
	node = &ast.FuncDecl{}
	if node.FunctionToken, err = p.Expect(token.FUNCTION); err != nil {
		return
	}
	if node.Name, err = p.Expect(token.IDENT); err != nil {
		return
	}
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return
	}
	body, err := ParseBlock(p)
	if err != nil {
		return
	}
	node.Body = body
	return
}

func ParseBlock(p *Parser) (node *ast.Block, err error) {
	p.EnterScope(blockScope)
	defer p.ExitScope(blockScope)
	node = &ast.Block{}
	if node.LbraceToken, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	var errs []error
	for p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE {
		stmt, err := p.ParseStmt()
		if err != nil {
			errs = append(errs, err)
			p.AdvanceToStmtEnd()
			continue
		}
		node.Stmts = append(node.Stmts, stmt)
	}
	if node.RbraceToken, err = p.Expect(token.RBRACE); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		err = errors.Join(errs...)
		return
	}
	return
}
