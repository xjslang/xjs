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

func ParseRightExpr(p *Parser) (val ast.Node, err error) {
	typ0 := p.CurrentToken.Type
	p.AdvanceToken()
	if val, err = ParseValue(p); err != nil {
		return
	}
	for {
		typ1 := p.CurrentToken.Type
		if !typ1.IsBinaryOp() || typ0.Precedence() >= typ1.Precedence() {
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
	if typ.IsUnaryOp() {
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

func ParseCallExpr(p *Parser, leftVal ast.Node) (node *ast.CallExpr, err error) {
	node = &ast.CallExpr{Function: leftVal}
	if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
		return
	}
	if p.CurrentToken.Type != token.RPAREN {
		for {
			val, err := p.ParseExpr()
			if err != nil {
				return nil, err
			}
			node.Arguments = append(node.Arguments, val)
			if p.CurrentToken.Type == token.RPAREN {
				break
			}
			if _, err := p.Expect(token.COMMA); err != nil {
				return nil, err
			}
		}
	}
	if node.RparenToken, err = p.Expect(token.RPAREN); err != nil {
		return nil, err
	}
	return
}
