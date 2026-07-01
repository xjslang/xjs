package parser

import (
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Scanner interface {
	NextToken() token.Token
}

type Range struct {
	Start token.Position `json:"start"`
	End   token.Position `json:"end"`
}

type Error struct {
	Range   Range  `json:"range"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	start := err.Range.Start
	return "[line:" + strconv.Itoa(start.Line) +
		", col:" + strconv.Itoa(start.Column) +
		"] " + err.Message
}

type ErrorList []error

func (list ErrorList) Error() string {
	s := strings.Builder{}
	for i, err := range list {
		if err == nil {
			continue
		}
		if i > 0 {
			s.WriteRune('\n')
		}
		s.WriteString(err.Error())
	}
	return s.String()
}

type Parser struct {
	CurrentToken     token.Token
	PeekToken        token.Token
	scanner          Scanner
	scopes           ScopeTracker
	stmtParser       func(p *Parser) (ast.Stmt, error)
	exprParser       func(p *Parser) (ast.Expr, error)
	binaryExprParser func(p *Parser, left ast.Expr) (ast.Expr, error)
	unaryExprParser  func(p *Parser) (ast.Expr, error)
}

func (p *Parser) init(sc Scanner) {
	p.scopes = make(ScopeTracker)
	p.scanner = sc
	if p.stmtParser == nil {
		p.stmtParser = func(p *Parser) (ast.Stmt, error) {
			return nil, errors.New("unknown statement")
		}
	}
	if p.exprParser == nil {
		p.exprParser = func(p *Parser) (ast.Expr, error) {
			return nil, errors.New("unknown expression")
		}
	}
	if p.binaryExprParser == nil {
		p.binaryExprParser = func(p *Parser, left ast.Expr) (ast.Expr, error) {
			return nil, errors.New("unknown binary expression")
		}
	}
	if p.unaryExprParser == nil {
		p.unaryExprParser = func(p *Parser) (ast.Expr, error) {
			return nil, errors.New("unknown unary expression")
		}
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

func (p *Parser) ParseBinaryExpr(left ast.Expr) (ast.Expr, error) {
	return p.binaryExprParser(p, left)
}

func (p *Parser) ParseUnaryExpr() (ast.Expr, error) {
	return p.unaryExprParser(p)
}

func (p *Parser) AdvanceToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.scanner.NextToken()
}

func (p *Parser) Expect(typ token.Type) (token.Token, error) {
	tok := p.CurrentToken
	if p.CurrentToken.Type != typ {
		return tok, p.Error(typ.String() + " expected")
	}
	p.AdvanceToken()
	return tok, nil
}

func (p *Parser) ExpectString(s string) (token.Token, error) {
	tok := p.CurrentToken
	if tok.Literal != s {
		return tok, p.Error(s + " expected")
	}
	p.AdvanceToken()
	return tok, nil
}

func (p *Parser) Error(msg string) error {
	return p.ErrorAt(p.CurrentToken, msg)
}

func (p *Parser) ErrorAt(tok token.Token, msg string) error {
	line := tok.Line
	column := tok.Column
	if tok.Type == token.EOF {
		column++
	}
	return Error{
		Range: Range{
			Start: token.Position{
				Line:   line,
				Column: column,
			},
			End: token.Position{
				Line:   line,
				Column: column + utf8.RuneCountInString(tok.Literal),
			},
		},
		Message: msg,
	}
}

func (p *Parser) EnterScope(sc Scope) {
	p.scopes.Enter(sc)
}

func (p *Parser) ExitScope(sc Scope) {
	p.scopes.Exit(sc)
}

func (p *Parser) InScope(sc Scope) bool {
	return p.scopes.In(sc)
}
