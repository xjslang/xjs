package parser

import (
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

func (err Error) Error() string {
	return err.Message
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
	lexer          *lexer.Lexer
	CurrentToken   token.Token
	PeekToken      token.Token
	infixOperators map[token.TokenType]infixOperator

	statementParser func(p *Parser) ast.Statement

	errors ErrorList
}

func (p *Parser) Init(l *lexer.Lexer) {
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
	result := p.parseBody()
	if len(p.errors) > 0 {
		return result, p.errors
	}
	return result, nil
}

func (p *Parser) ParseExpression() (ast.Expression, error) {
	result := p.parseExpression()
	if len(p.errors) > 0 {
		return result, p.errors
	}
	return result, nil
}

func (p *Parser) parseExpression() ast.Expression {
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
	parseTerm := func() (ast.Expression, token.Token) {
		// parse val
		val := p.parseValue()
		if val == nil {
			return nil, token.Token{}
		}
		// parse op
		op := p.CurrentToken
		if registered(op.Type) {
			p.AdvanceToken()
		}
		return val, op
	}
	var parseRightExp func(ast.Expression, token.Token) (ast.Expression, token.Token)
	parseRightExp = func(v0 ast.Expression, op0 token.Token) (ast.Expression, token.Token) {
		for {
			v1, op1 := parseTerm()
			if v1 == nil {
				return nil, token.Token{}
			}
			if precedence(op0.Type) < precedence(op1.Type) {
				v1, op1 = parseRightExp(v1, op1)
				if v1 == nil {
					return nil, token.Token{}
				}
			}
			v0 = p.infixOperators[op0.Type].fn(op0, v0, v1)
			if precedence(op0.Type) > precedence(op1.Type) {
				return v0, op1
			}
			op0 = op1
		}
	}

	v, op := parseTerm()
	if v == nil {
		return nil
	}
	for registered(op.Type) {
		v, op = parseRightExp(v, op)
		if v == nil {
			return nil
		}
	}
	return v
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

// got checks that the current token matches the expected type,
// advances the position, and returns the token.
//
// If the token does not match, it records an error and returns it.
func (p *Parser) got(ttype token.TokenType) bool {
	if p.CurrentToken.Type != ttype {
		p.AddError("Expected " + ttype.String())
		return false
	}
	p.AdvanceToken()
	return true
}

func (p *Parser) gotStatementTerminator() bool {
	if p.CurrentToken.Type == token.SEMICOLON {
		p.AdvanceToken()
		return true
	}
	if p.CurrentToken.Type == token.EOF || p.CurrentToken.AfterNewline {
		return true
	}
	p.AddError("Expected statement terminator")
	return false
}

func (p *Parser) advanceToStatementEnd() {
	for {
		if p.CurrentToken.Type == token.SEMICOLON {
			p.AdvanceToken()
			break
		}
		if p.CurrentToken.Type == token.EOF || p.CurrentToken.AfterNewline {
			break
		}
		p.AdvanceToken()
	}
}

func (p *Parser) parseBody() *ast.BlockStatement {
	bodyStmt := &ast.BlockStatement{}
	for stmt := p.statementParser(p); stmt != nil; stmt = p.statementParser(p) {
		p.advanceToStatementEnd()
		bodyStmt.Statements = append(bodyStmt.Statements, stmt)
	}
	return bodyStmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{}
	p.AdvanceToken() // consume token.LET
	ident := p.CurrentToken
	if !p.got(token.IDENT) {
		return nil
	}
	stmt.Name = ident
	if !p.got(token.ASSIGN) {
		return nil
	}
	stmt.Value = p.parseExpression()
	if stmt.Value == nil {
		return nil
	}
	if !p.gotStatementTerminator() {
		return nil
	}
	return stmt
}

func (p *Parser) parseFunction() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{}
	p.AdvanceToken() // consume token.FUNCTION
	ident := p.CurrentToken
	if !p.got(token.IDENT) {
		return nil
	}
	stmt.Name = ident
	if !p.got(token.LPAREN) {
		return nil
	}
	if !p.got(token.RPAREN) {
		return nil
	}
	if !p.got(token.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBody()
	if !p.got(token.RBRACE) {
		return nil
	}
	return stmt
}

func (p *Parser) parseValue() ast.Expression {
	switch p.CurrentToken.Type {
	case token.LPAREN:
		p.AdvanceToken() // consume (
		exp := p.parseExpression()
		if exp == nil {
			return nil
		}
		if !p.got(token.RPAREN) {
			return nil
		}
		return &ast.GroupedExpression{Value: exp}
	case token.NUMBER:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.IntegerLiteral{Value: val}
	case token.STRING:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.StringLiteral{Value: val}
	case token.BOOLEAN:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.BooleanLiteral{Value: val}
	case token.IDENT:
		val := p.CurrentToken.Literal
		p.AdvanceToken()
		return &ast.Identifier{Value: val}
	}
	p.AddError("Expected value")
	return nil
}
