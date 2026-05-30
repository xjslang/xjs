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
		prevToken := p.CurrentToken
		stmt, err := p.ParseStmt()
		if err != nil {
			if prevToken.Position == p.CurrentToken.Position {
				// advance position to avoid infinite loop
				p.AdvanceToken()
			}
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

func ParseBlock(p *Parser) (node *ast.Block, err error) {
	p.EnterScope(blockScope)
	defer p.ExitScope(blockScope)
	node = &ast.Block{}
	if node.LbraceToken, err = p.Expect(token.LBRACE); err != nil {
		return
	}
	var errs []error
	for p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE {
		prevToken := p.CurrentToken
		stmt, err := p.ParseStmt()
		if err != nil {
			if prevToken.Position == p.CurrentToken.Position {
				// advance position to avoid infinite loop
				p.AdvanceToken()
			}
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

func ParseRightValue(p *Parser, precedence int) (val ast.Node, err error) {
	if val, err = ParseValue(p); err != nil {
		return
	}
	for {
		typ := p.CurrentToken.Type
		if !typ.IsInfixOp() || precedence >= typ.Precedence() {
			break
		}
		if val, err = p.infixExprParser(p, val); err != nil {
			return
		}
	}
	return
}

func ParseValue(p *Parser) (ast.Node, error) {
	typ := p.CurrentToken.Type
	if typ.IsPrefixOp() {
		return p.prefixExprParser(p)
	}
	switch typ {
	case token.NUMBER, token.STRING, token.BOOLEAN:
		val := p.CurrentToken
		p.AdvanceToken()
		return &ast.BasicLit{Value: val}, nil
	case token.IDENT:
		val := p.CurrentToken
		p.AdvanceToken()
		return &ast.Ident{Value: val}, nil
	}
	msg := "Expected value"
	p.AddError(msg)
	return nil, errors.New(msg)
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
