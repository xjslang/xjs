package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Parser struct {
	CurrentToken   token.Token
	PeekToken      token.Token
	scanner        *Scanner
	stmtParser     func(p *Parser) (ast.Stmt, error)
	exprParser     func(p *Parser) (ast.Expr, error)
	infixOpParser  func(p *Parser, left ast.Expr) (ast.Expr, error)
	prefixOpParser func(p *Parser) (ast.Expr, error)
}

func (p *Parser) init(s *Scanner) {
	p.scanner = s
	if p.stmtParser == nil {
		p.stmtParser = defaultStmtParser
	}
	if p.exprParser == nil {
		p.exprParser = defaultExprParser
	}
	if p.infixOpParser == nil {
		p.infixOpParser = infixOpParser
	}
	if p.prefixOpParser == nil {
		p.prefixOpParser = defaultPrefixOpParser
	}
	p.CurrentToken = token.Token{}
	p.PeekToken = token.Token{}
	// call twice to update CurrentToken and PeekToken
	p.AdvanceToken()
	p.AdvanceToken()
}

func (p *Parser) ParseStmt() (ast.Stmt, error) {
	return p.stmtParser(p)
}

func (p *Parser) ParseExpr() (ast.Expr, error) {
	return p.exprParser(p)
}

func (p *Parser) ParseInfixOp(left ast.Expr) (ast.Expr, error) {
	return p.infixOpParser(p, left)
}

func (p *Parser) ParsePrefixOp() (ast.Expr, error) {
	return p.prefixOpParser(p)
}

func (p *Parser) AdvanceToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.scanner.NextToken()
}

// Expect consumes the current token if it matches the expected type, or returns an error.
func (p *Parser) Expect(typ token.Type) (token.Token, error) {
	tok := p.CurrentToken
	if p.CurrentToken.Type != typ {
		s := typ.Literal()
		if s == "" {
			s = typ.Name()
		}
		return tok, p.Error(s + " expected")
	}
	p.AdvanceToken()
	return tok, nil
}

// ExpectLiteral consumes the current token if it matches the expected literal, or returns an error.
func (p *Parser) ExpectLiteral(s string) (token.Token, error) {
	tok := p.CurrentToken
	if tok.Literal != s {
		return tok, p.Error(s + " expected")
	}
	p.AdvanceToken()
	return tok, nil
}

// ExpectWith consumes the current token if it is accepted by the given scanner function, or returns an error.
func (p *Parser) ExpectWith(scanner func(*Scanner) (string, error)) (tok token.Token, err error) {
	tok = p.CurrentToken
	f := p.scanner.ForkFrom(p.CurrentToken.Offset)
	if tok.Literal, err = scanner(f); err != nil {
		return
	}
	p.scanner.Apply(f)
	p.AdvanceToken()
	p.CurrentToken = tok
	return
}

// ExpectSemi expects a semicolon or any other symbol that acts as a
// "statement terminator", such as '}' or ')'. If the statement terminator is a
// semicolon, then it consumes it and advances to the next token.
func (p *Parser) ExpectSemi() (tok token.Token, err error) {
	tok = p.CurrentToken
	switch tok.Type {
	case token.SEMICOLON:
		p.AdvanceToken()
		return
	case token.RBRACE, token.RPAREN, token.EOF:
		tok = token.Token{
			Type:    token.SEMICOLON,
			Literal: token.SEMICOLON.Literal(),
			Offset:  tok.Offset,
		}
		return
	default:
		if tok.IsAfterNewline() {
			tok = token.Token{
				Type:    token.SEMICOLON,
				Literal: token.SEMICOLON.Literal(),
				Offset:  tok.Offset,
			}
			return
		}
	}
	err = p.Error(token.SEMICOLON.Literal() + " expected")
	return
}

// Error returns an error at the current token position.
func (p *Parser) Error(msg string) error {
	return p.ErrorAt(p.CurrentToken, msg)
}

// ErrorAt returns an error at the given token position.
func (p *Parser) ErrorAt(tok token.Token, msg string) error {
	offset := tok.Offset
	if tok.Type == token.EOF {
		offset++
	}
	return token.Error{
		Range:   token.Range{offset, offset + len(tok.Literal)},
		Message: msg,
	}
}
