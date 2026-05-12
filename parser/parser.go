package parser

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
)

type Range struct {
	Start scanner.Position `json:"start"`
	End   scanner.Position `json:"end"`
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

type infixOperator struct {
	precedence int
	fn         func(op scanner.Token, left, right ast.Node) ast.Node
}

type Parser struct {
	scanner        *scanner.Scanner
	CurrentToken   scanner.Token
	PeekToken      scanner.Token
	scopes         ScopeTracker
	infixOperators map[scanner.Kind]infixOperator

	statementParser func(p *Parser) (ast.Node, error)

	errors ErrorList
}

func (p *Parser) Init(sc *scanner.Scanner) {
	p.scopes = make(ScopeTracker)
	p.scanner = sc
	if p.statementParser == nil {
		p.statementParser = defaultStatementParser
	}
	p.CurrentToken = scanner.Token{}
	p.PeekToken = scanner.Token{}
	if p.infixOperators == nil {
		p.infixOperators = make(map[scanner.Kind]infixOperator)
	}
	p.infixOperators[scanner.PLUS] = infixOperator{precedence: 1, fn: defaultInfixOperator}
	p.infixOperators[scanner.MINUS] = infixOperator{precedence: 1, fn: defaultInfixOperator}
	p.infixOperators[scanner.MULTIPLY] = infixOperator{precedence: 2, fn: defaultInfixOperator}
	p.infixOperators[scanner.DIVIDE] = infixOperator{precedence: 2, fn: defaultInfixOperator}
	p.infixOperators[scanner.MODULO] = infixOperator{precedence: 2, fn: defaultInfixOperator}
	p.errors = ErrorList{}
	// call twice to update CurrentToken and PeekToken
	p.AdvanceToken()
	p.AdvanceToken()
}

func defaultInfixOperator(op scanner.Token, left, right ast.Node) ast.Node {
	return &ast.InfixOperator{
		LeftValue:  left,
		Operator:   op,
		RightValue: right,
	}
}

func (p *Parser) ParseStatement() (ast.Node, error) {
	return p.statementParser(p)
}

func (p *Parser) ParseExpression() (ast.Node, error) {
	registered := func(tt scanner.Kind) bool {
		_, ok := p.infixOperators[tt]
		return ok
	}
	precedence := func(tt scanner.Kind) int {
		if op, ok := p.infixOperators[tt]; ok {
			return op.precedence
		}
		return -1
	}
	parseTerm := func() (ast.Node, scanner.Token, error) {
		// parse val
		val, err := p.parseValue()
		if err != nil {
			return nil, scanner.Token{}, err
		}
		// parse op
		op := p.CurrentToken
		if registered(op.Type) {
			p.AdvanceToken()
		}
		return val, op, nil
	}
	var parseRightExp func(ast.Node, scanner.Token) (ast.Node, scanner.Token, error)
	parseRightExp = func(v0 ast.Node, op0 scanner.Token) (ast.Node, scanner.Token, error) {
		for {
			v1, op1, err := parseTerm()
			if err != nil {
				return nil, scanner.Token{}, err
			}
			if precedence(op0.Type) < precedence(op1.Type) {
				v1, op1, err = parseRightExp(v1, op1)
				if err != nil {
					return nil, scanner.Token{}, err
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
		Range: Range{
			Start: scanner.Position{
				Line:   line,
				Column: column,
			},
			End: scanner.Position{
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

// Expect checks that the current token matches the expected type and advances the position.
//
// If the token does not match, it records an error and returns it.
func (p *Parser) Expect(ttype scanner.Kind) error {
	if p.CurrentToken.Type != ttype {
		msg := "Expected " + ttype.String()
		p.AddError(msg)
		return errors.New(msg)
	}
	p.AdvanceToken()
	return nil
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

func (p *Parser) parseValue() (ast.Node, error) {
	switch p.CurrentToken.Type {
	case scanner.LPAREN:
		p.AdvanceToken() // consume (
		exp, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.Expect(scanner.RPAREN); err != nil {
			return nil, err
		}
		return &ast.GroupedExpression{Value: exp}, nil
	case scanner.NUMBER:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.Integer{Value: val}, nil
	case scanner.STRING:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.String{Value: val}, nil
	case scanner.BOOLEAN:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.Boolean{Value: val}, nil
	case scanner.IDENT:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.Ident{Value: val}, nil
	}
	msg := "Expected value"
	p.AddError(msg)
	return nil, errors.New(msg)
}
