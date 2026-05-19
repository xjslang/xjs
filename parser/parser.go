package parser

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

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

type infixOperator struct {
	precedence int
	fn         func(op token.Token, left, right ast.Node) ast.Node
}

type Parser struct {
	CurrentToken    token.Token
	PeekToken       token.Token
	scanner         *scanner.Scanner
	scopes          ScopeTracker
	infixOperators  map[token.Type]infixOperator
	statementParser func(p *Parser) (ast.Node, error)
	errors          ErrorList
}

func (p *Parser) Init(sc *scanner.Scanner) {
	p.scopes = make(ScopeTracker)
	p.scanner = sc
	if p.statementParser == nil {
		p.statementParser = defaultStatementParser
	}
	p.CurrentToken = token.Token{}
	p.PeekToken = token.Token{}
	if p.infixOperators == nil {
		p.infixOperators = make(map[token.Type]infixOperator)
	}
	defaultInfixOperator := func(op token.Token, left, right ast.Node) ast.Node {
		return &ast.BinaryExpr{
			LeftValue:  left,
			Operator:   op,
			RightValue: right,
		}
	}
	p.infixOperators[token.PLUS] = infixOperator{precedence: 1, fn: defaultInfixOperator}
	p.infixOperators[token.MINUS] = infixOperator{precedence: 1, fn: defaultInfixOperator}
	p.infixOperators[token.MULTIPLY] = infixOperator{precedence: 2, fn: defaultInfixOperator}
	p.infixOperators[token.DIVIDE] = infixOperator{precedence: 2, fn: defaultInfixOperator}
	p.infixOperators[token.MODULO] = infixOperator{precedence: 2, fn: defaultInfixOperator}
	p.infixOperators[token.LPAREN] = infixOperator{precedence: 3, fn: defaultInfixOperator}
	p.errors = ErrorList{}
	// call twice to update CurrentToken and PeekToken
	p.AdvanceToken()
	p.AdvanceToken()
}

func (p *Parser) ParseStatement() (ast.Node, error) {
	return p.statementParser(p)
}

func (p *Parser) ParseExpression() (ast.Node, error) {
	registered := func(typ token.Type) bool {
		_, ok := p.infixOperators[typ]
		return ok
	}
	precedence := func(typ token.Type) int {
		if op, ok := p.infixOperators[typ]; ok {
			return op.precedence
		}
		return -1
	}
	parseTerm := func() (ast.Node, token.Token, error) {
		// parse val
		val, err := p.parseValue()
		if err != nil {
			return nil, token.Token{}, err
		}
		// parse op
		op := p.CurrentToken
		return val, op, nil
	}
	parseInfixCall := func(val ast.Node) (node *ast.CallExpr, err error) {
		node = &ast.CallExpr{
			Function:  val,
			Arguments: nil,
		}
		if node.LparenToken, err = p.Expect(token.LPAREN); err != nil {
			return nil, err
		}
		if p.CurrentToken.Type != token.RPAREN {
			for {
				val, err := p.ParseExpression()
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
		return node, nil
	}
	var parseRightExp func(ast.Node, token.Token) (ast.Node, error)
	parseRightExp = func(v0 ast.Node, op0 token.Token) (node ast.Node, err error) {
		if op0.Type == token.LPAREN {
			return parseInfixCall(v0)
		}
		for {
			p.AdvanceToken()
			v1, op1, err := parseTerm()
			if err != nil {
				return nil, err
			}
			if precedence(op0.Type) < precedence(op1.Type) {
				v1, err = parseRightExp(v1, op1)
				if err != nil {
					return nil, err
				}
				op1 = p.CurrentToken
			}
			v0 = p.infixOperators[op0.Type].fn(op0, v0, v1)
			if precedence(op0.Type) > precedence(op1.Type) {
				return v0, nil
			}
			op0 = op1
		}
	}
	v, op, err := parseTerm()
	if err != nil {
		return nil, err
	}
	for registered(op.Type) {
		v, err = parseRightExp(v, op)
		if err != nil {
			return nil, err
		}
		op = p.CurrentToken
	}
	return v, nil
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
	if tok.Type == token.EOF || tok.AfterNewline {
		tok = token.Token{Type: token.SEMICOLON, Literal: token.SEMICOLON.String(), Position: tok.Position}
		return tok, nil
	}
	if p.InScope(blockScope) && tok.Type == token.RBRACE {
		tok = token.Token{Type: token.SEMICOLON, Literal: token.SEMICOLON.String(), Position: tok.Position}
		return tok, nil
	}
	msg := "Expected statement terminator"
	p.AddError(msg)
	return tok, errors.New(msg)
}

func (p *Parser) AdvanceToStatementEnd() {
	for {
		typ := p.CurrentToken.Type
		if typ == token.SEMICOLON {
			p.AdvanceToken()
			break
		}
		if typ == token.EOF || p.CurrentToken.AfterNewline ||
			p.InScope(blockScope) && typ == token.RBRACE {
			break
		}
		p.AdvanceToken()
	}
}

func (p *Parser) parseValue() (ast.Node, error) {
	switch p.CurrentToken.Type {
	case token.LPAREN:
		return ParseParenExpr(p)
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
