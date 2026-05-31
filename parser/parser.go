package parser

import (
	"errors"
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

type ErrorList []Error

func (list ErrorList) Error() string {
	var result []string
	for _, err := range list {
		result = append(result, err.Message)
	}
	return strings.Join(result, "\n")
}

type Parser struct {
	CurrentToken     token.Token
	PeekToken        token.Token
	scanner          Scanner
	scopes           ScopeTracker
	stmtParser       func(p *Parser) (ast.Node, error)
	exprParser       func(p *Parser) (ast.Node, error)
	binaryExprParser func(p *Parser, leftVal ast.Node) (ast.Node, error)
	unaryExprParser  func(p *Parser) (ast.Node, error)
	errors           ErrorList
}

// Init initializes the parser.
//
// Call Init before parsing with Parse/ParseExprStmt/ParseExpr.
// Parser middleware must be registered via Use*Parser methods BEFORE Init.
func (p *Parser) Init(sc Scanner) {
	p.scopes = make(ScopeTracker)
	p.scanner = sc
	if p.stmtParser == nil {
		p.stmtParser = func(p *Parser) (ast.Node, error) {
			return nil, errors.New("unknown statement")
		}
	}
	if p.exprParser == nil {
		p.exprParser = func(p *Parser) (ast.Node, error) {
			return nil, errors.New("unknown expression")
		}
	}
	if p.binaryExprParser == nil {
		p.binaryExprParser = func(p *Parser, leftVal ast.Node) (ast.Node, error) {
			return nil, errors.New("unknown binary expression")
		}
	}
	if p.unaryExprParser == nil {
		p.unaryExprParser = func(p *Parser) (ast.Node, error) {
			return nil, errors.New("unknown unary expression")
		}
	}
	p.CurrentToken = token.Token{}
	p.PeekToken = token.Token{}
	p.errors = ErrorList{}
	// call twice to update CurrentToken and PeekToken
	p.AdvanceToken()
	p.AdvanceToken()
}

func (p *Parser) ParseExprStmt() (ast.Node, error) {
	return p.stmtParser(p)
}

func (p *Parser) ParseExpr() (ast.Node, error) {
	return p.exprParser(p)
}

func (p *Parser) ParseBinaryExpr(leftVal ast.Node) (ast.Node, error) {
	return p.binaryExprParser(p, leftVal)
}

func (p *Parser) ParseUnaryExpr() (ast.Node, error) {
	return p.unaryExprParser(p)
}

func (p *Parser) AddError(msg string) {
	line := p.CurrentToken.Line
	column := p.CurrentToken.Column
	p.errors = append(p.errors, Error{
		Range: Range{
			Start: token.Position{
				Line:   line,
				Column: column,
			},
			End: token.Position{
				Line:   line,
				Column: column + utf8.RuneCountInString(p.CurrentToken.Literal),
			},
		},
		Message: msg,
	})
}

func (p *Parser) AdvanceToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.scanner.NextToken()
}

func (p *Parser) Expect(typ token.Type) (token.Token, error) {
	tok := p.CurrentToken
	if p.CurrentToken.Type != typ {
		msg := "Expected " + typ.String()
		p.AddError(msg)
		return tok, errors.New(msg)
	}
	p.AdvanceToken()
	return tok, nil
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

func (p *Parser) Errors() ErrorList {
	return append(ErrorList{}, p.errors...)
}

func (p *Parser) ExpectSemi() (token.Token, error) {
	tok := p.CurrentToken
	if tok.Type == token.SEMICOLON {
		p.AdvanceToken()
		return tok, nil
	}
	if tok.Type == token.EOF || tok.AfterNewline || tok.Type == token.RBRACE {
		tok = token.Token{Type: token.SEMICOLON, Literal: token.SEMICOLON.String(), Position: tok.Position}
		return tok, nil
	}
	msg := "Expected statement terminator"
	p.AddError(msg)
	return tok, errors.New(msg)
}

func (p *Parser) AdvanceToStmtEnd() {
	for {
		typ := p.CurrentToken.Type
		if typ == token.SEMICOLON {
			p.AdvanceToken()
			break
		}
		if typ == token.EOF || p.CurrentToken.AfterNewline || typ == token.RBRACE {
			break
		}
		p.AdvanceToken()
	}
}
