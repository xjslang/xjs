package parser

import (
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

type Range struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type Error struct {
	Range   Range  `json:"range"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
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
	CurrentToken   token.Token
	PeekToken      token.Token
	scanner        token.Scanner
	stmtParser     func(p *Parser) (ast.Stmt, error)
	exprParser     func(p *Parser) (ast.Expr, error)
	infixOpParser  func(p *Parser, left ast.Expr) (ast.Expr, error)
	prefixOpParser func(p *Parser) (ast.Expr, error)
}

func (p *Parser) init(s token.Scanner) {
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

// Fork returns a copy of the parser that can advance without modifying the original state.
//
// See also: [Parser.Apply]
func (p *Parser) Fork() *Parser {
	s := p.scanner.(token.ForkableScanner)
	return &Parser{
		CurrentToken:   p.CurrentToken,
		PeekToken:      p.PeekToken,
		scanner:        s.Fork(),
		stmtParser:     p.stmtParser,
		exprParser:     p.exprParser,
		infixOpParser:  p.infixOpParser,
		prefixOpParser: p.prefixOpParser,
	}
}

// Apply copies the state from a forked parser back into the original parser.
//
// See also: [Parser.Fork]
func (p *Parser) Apply(forked *Parser) {
	s := p.scanner.(token.ForkableScanner)
	s.Apply(forked.scanner)
	p.CurrentToken = forked.CurrentToken
	p.PeekToken = forked.PeekToken
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
func (p *Parser) ExpectWith(scanner func(s token.Scanner) (string, error)) (tok token.Token, err error) {
	tok = p.CurrentToken
	s := p.scanner.(token.ForkableScanner)
	f := s.ForkFrom(p.CurrentToken.Offset)
	if tok.Literal, err = scanner(f); err != nil {
		return
	}
	s.Apply(f)
	p.AdvanceToken()
	p.CurrentToken = tok
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
	return Error{
		Range: Range{
			Start: offset,
			End:   offset + len(tok.Literal),
		},
		Message: msg,
	}
}
