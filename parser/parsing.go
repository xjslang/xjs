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

func ParseRemainingExpr(p *Parser) (val ast.Node, err error) {
	acc := func(leftVal, rightVal ast.Node, op token.Token) ast.Node {
		if leftVal == nil {
			return rightVal
		}
		return &ast.BinaryExpr{
			LeftValue:  leftVal,
			Operator:   op,
			RightValue: rightVal,
		}
	}
	var nextVal ast.Node
	op := p.CurrentToken
	p.AdvanceToken()
	for {
		if nextVal, err = p.parseValue(); err != nil {
			return
		}
		nextOp := p.CurrentToken
		if !p.isOperator(nextOp) {
			val = acc(val, nextVal, op)
			break
		}
		p0, p1 := p.precedence(op), p.precedence(nextOp)
		if p0 < p1 {
			var remainingVal ast.Node
			if remainingVal, err = ParseRemainingExpr(p); err != nil {
				return
			}
			nextVal = acc(nextVal, remainingVal, nextOp)
			nextOp = p.CurrentToken
		} else if p0 > p1 {
			val = acc(val, nextVal, op)
			break
		}
		p.AdvanceToken()
		val = acc(val, nextVal, op)
		op = nextOp
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
