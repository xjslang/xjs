package parser

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/source"
	"github.com/xjslang/xjs/token"
)

var operatorPrecedence map[token.TokenType]int = map[token.TokenType]int{
	token.PLUS:     1,
	token.MINUS:    1,
	token.MULTIPLY: 2,
	token.DIVIDE:   2,
	token.MODULO:   2,
}

type Error struct {
	Range   source.Range `json:"range"`
	Message string       `json:"message"`
}

type ErrorList []Error

func (list ErrorList) Error() string {
	var result []string
	for _, err := range list {
		result = append(result, err.Message)
	}
	return strings.Join(result, "\n")
}

type infixOperator struct {
	precedence int
	fn         func(op token.Token, left, right ast.Expression) ast.Expression
}

type Parser struct {
	scopes         scopeTracker
	lexer          *lexer.Lexer
	CurrentToken   token.Token
	PeekToken      token.Token
	infixOperators map[token.TokenType]infixOperator

	statementParser func(p *Parser) (ast.Statement, error)

	errors ErrorList
}

func (p *Parser) Init(l *lexer.Lexer) {
	p.scopes = make(scopeTracker)
	p.lexer = l
	if p.statementParser == nil {
		p.statementParser = defaultStatementParser
	}
	p.CurrentToken = token.Token{}
	p.PeekToken = token.Token{}
	if p.infixOperators == nil {
		p.infixOperators = make(map[token.TokenType]infixOperator)
	}
	for op, precedence := range operatorPrecedence {
		if _, ok := p.infixOperators[op]; ok {
			continue
		}
		p.infixOperators[op] = infixOperator{
			precedence: precedence,
			fn:         defaultInfixOperator,
		}
	}
	p.errors = ErrorList{}
	// call twice to update CurrentToken and PeekToken
	p.AdvanceToken()
	p.AdvanceToken()
}

func (p *Parser) ParseStatement() (ast.Statement, error) {
	return p.statementParser(p)
}

func (p *Parser) ParseExpression() (ast.Expression, error) {
	registered := func(tt token.TokenType) bool {
		_, ok := p.infixOperators[tt]
		return ok
	}
	precedence := func(tt token.TokenType) int {
		if op, ok := p.infixOperators[tt]; ok {
			return op.precedence
		}
		return -1
	}
	parseTerm := func() (ast.Expression, token.Token, error) {
		// parse val
		val, err := p.parseValue()
		if err != nil {
			return nil, token.Token{}, err
		}
		// parse op
		op := p.CurrentToken
		if registered(op.Type) {
			p.AdvanceToken()
		}
		return val, op, nil
	}
	var parseRightExp func(ast.Expression, token.Token) (ast.Expression, token.Token, error)
	parseRightExp = func(v0 ast.Expression, op0 token.Token) (ast.Expression, token.Token, error) {
		for {
			v1, op1, err := parseTerm()
			if err != nil {
				return nil, token.Token{}, err
			}
			if precedence(op0.Type) < precedence(op1.Type) {
				v1, op1, err = parseRightExp(v1, op1)
				if err != nil {
					return nil, token.Token{}, err
				}
			}
			v0 = p.infixOperators[op0.Type].fn(op0, v0, v1)
			if precedence(op0.Type) > precedence(op1.Type) {
				return v0, op1, nil
			}
			op0 = op1
		}
	}

	v, op, err := parseTerm()
	if err != nil {
		return nil, err
	}
	for registered(op.Type) {
		v, op, err = parseRightExp(v, op)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func (p *Parser) AddError(msg string) {
	line := p.CurrentToken.Line
	column := p.CurrentToken.Column
	p.errors = append(p.errors, Error{
		Range: source.Range{
			Start: source.Position{
				Line:   line,
				Column: column,
			},
			End: source.Position{
				Line:   line,
				Column: column + utf8.RuneCountInString(p.CurrentToken.Literal),
			},
		},
		Message: msg,
	})
}

func (p *Parser) AdvanceToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

// Expect checks that the current token matches the expected type and advances the position.
//
// If the token does not match, it records an error and returns it.
func (p *Parser) Expect(ttype token.TokenType) error {
	if p.CurrentToken.Type != ttype {
		msg := "Expected " + ttype.String()
		p.AddError(msg)
		return errors.New(msg)
	}
	p.AdvanceToken()
	return nil
}

func (p *Parser) EnterScope(sc scope) {
	p.scopes.Enter(sc)
}

func (p *Parser) ExitScope(sc scope) {
	p.scopes.Exit(sc)
}

func (p *Parser) InScope(sc scope) bool {
	return p.scopes.In(sc)
}

func (p *Parser) Errors() ErrorList {
	return append(ErrorList{}, p.errors...)
}

func (p *Parser) parseValue() (ast.Expression, error) {
	switch p.CurrentToken.Type {
	case token.LPAREN:
		p.AdvanceToken() // consume (
		exp, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.Expect(token.RPAREN); err != nil {
			return nil, err
		}
		return &ast.GroupedExpression{Value: exp}, nil
	case token.NUMBER:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.IntegerLiteral{Value: val}, nil
	case token.STRING:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.StringLiteral{Value: val}, nil
	case token.BOOLEAN:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.BooleanLiteral{Value: val}, nil
	case token.IDENT:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.Identifier{Value: val}, nil
	}
	msg := "Expected value"
	p.AddError(msg)
	return nil, errors.New(msg)
}
