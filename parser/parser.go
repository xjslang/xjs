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

func Parse(input []byte) (*ast.BlockStatement, error) {
	l := &lexer.Lexer{}
	l.Init(input)
	p := Parser{}
	p.Init(l)
	return p.ParseProgram()
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

func (p *Parser) ParseProgram() (*ast.BlockStatement, error) {
	result := &ast.BlockStatement{}
	for p.CurrentToken.Type != token.EOF {
		stmt, err := p.statementParser(p)
		if err != nil {
			p.advanceToStatementEnd()
			continue
		}
		result.Statements = append(result.Statements, stmt)
	}
	if len(p.errors) > 0 {
		return result, p.errors
	}
	return result, nil
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

// Expect checks that the current token matches the expected type,
// advances the position, and returns the token.
//
// If the token does not match, it records an error and returns it.
func (p *Parser) expect(ttype token.TokenType) (token.Token, error) {
	if p.CurrentToken.Type != ttype {
		msg := "Expected " + ttype.String()
		p.AddError(msg)
		return token.Token{}, errors.New(msg)
	}
	tok := p.CurrentToken
	p.AdvanceToken()
	return tok, nil
}

func (p *Parser) consumeStatementTerminator() bool {
	if p.CurrentToken.Type == token.SEMICOLON {
		p.AdvanceToken()
		return true
	}
	if p.CurrentToken.Type == token.EOF || p.CurrentToken.AfterNewline {
		return true
	}
	if p.scopes.in(blockScope) && p.CurrentToken.Type == token.RBRACE {
		return true
	}
	return false
}

func (p *Parser) expectStatementTerminator() error {
	if p.consumeStatementTerminator() {
		return nil
	}
	msg := "Expected statement terminator"
	p.AddError(msg)
	return errors.New(msg)
}

func (p *Parser) advanceToStatementEnd() {
	for !p.consumeStatementTerminator() {
		p.AdvanceToken()
	}
}

func (p *Parser) parseBody() *ast.BlockStatement {
	p.scopes.enter(blockScope)
	defer func() {
		p.scopes.exit(blockScope)
	}()
	bodyStmt := &ast.BlockStatement{}
	for p.CurrentToken.Type != token.EOF && p.CurrentToken.Type != token.RBRACE {
		stmt, err := p.statementParser(p)
		if err != nil {
			p.advanceToStatementEnd()
			continue
		}
		bodyStmt.Statements = append(bodyStmt.Statements, stmt)
	}
	return bodyStmt
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{}
	p.AdvanceToken() // consume token.LET
	ident, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}
	stmt.Name = ident
	if _, err := p.expect(token.ASSIGN); err != nil {
		return nil, err
	}
	stmt.Value, err = p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.expectStatementTerminator(); err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseFunction() (*ast.FunctionDeclaration, error) {
	stmt := &ast.FunctionDeclaration{}
	p.AdvanceToken() // consume token.FUNCTION
	ident, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}
	stmt.Name = ident
	if _, err := p.expect(token.LPAREN); err != nil {
		return nil, err
	}
	if _, err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	if _, err := p.expect(token.LBRACE); err != nil {
		return nil, err
	}
	stmt.Body = p.parseBody()
	if _, err := p.expect(token.RBRACE); err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseValue() (ast.Expression, error) {
	switch p.CurrentToken.Type {
	case token.LPAREN:
		p.AdvanceToken() // consume (
		exp, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RPAREN); err != nil {
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
