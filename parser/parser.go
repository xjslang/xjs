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

func Parse(input []byte) (*ast.BlockStatement, error) {
	l := lexer.New(input)
	p := newParser(l)
	return p.parseProgram()
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

type parser struct {
	lexer        *lexer.Lexer
	CurrentToken token.Token
	PeekToken    token.Token

	errors ErrorList
}

func newParser(l *lexer.Lexer) *parser {
	p := &parser{lexer: l}
	// call twice to update CurrentToken and PeekToken
	p.advanceToken()
	p.advanceToken()
	return p
}

func (p *parser) parseProgram() (*ast.BlockStatement, error) {
	result := p.parseBody()
	if len(p.errors) > 0 {
		return result, p.errors
	}
	return result, nil
}

func (p *parser) addError(msg string) {
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

func (p *parser) advanceToken() {
	p.CurrentToken = p.PeekToken
	p.PeekToken = p.lexer.NextToken()
}

// Expect checks that the current token matches the expected type,
// advances the position, and returns the token.
//
// If the token does not match, it records an error and returns it.
func (p *parser) expect(ttype token.TokenType) (token.Token, error) {
	if p.CurrentToken.Type != ttype {
		msg := "Expected " + ttype.String() + ", got " + p.CurrentToken.Type.String()
		p.addError(msg)
		return token.Token{}, errors.New(msg)
	}
	tok := p.CurrentToken
	p.advanceToken()
	return tok, nil
}

func (p *parser) expectStatementTerminator() error {
	if p.CurrentToken.Type == token.SEMICOLON {
		p.advanceToken()
		return nil
	}
	if p.CurrentToken.Type == token.EOF || p.CurrentToken.AfterNewline {
		return nil
	}
	msg := "Expected semicolon, newline, or EOF, got " + p.CurrentToken.Type.String()
	p.addError(msg)
	return errors.New(msg)
}

func (p *parser) parseBody() *ast.BlockStatement {
	bodyStmt := &ast.BlockStatement{}
	for stmt := p.parseStatement(); stmt != nil; stmt = p.parseStatement() {
		bodyStmt.Statements = append(bodyStmt.Statements, stmt)
	}
	return bodyStmt
}

func (p *parser) parseStatement() ast.Statement {
	switch p.CurrentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.FUNCTION:
		return p.parseFunction()
	}
	return nil
}

func (p *parser) parseExpression() ast.Expression {
	switch p.CurrentToken.Type {
	case token.NUMBER:
		val := p.CurrentToken.Literal
		p.advanceToken()
		return &ast.IntegerLiteral{Value: val}
	case token.STRING:
		val := p.CurrentToken.Literal
		p.advanceToken()
		return &ast.StringLiteral{Value: val}
	}
	return nil
}

func (p *parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{}
	p.advanceToken() // consume token.LET
	ident, err := p.expect(token.IDENT)
	if err != nil {
		return nil
	}
	stmt.Name = &ast.Identifier{Value: ident.Literal}
	if _, err := p.expect(token.ASSIGN); err != nil {
		return nil
	}
	stmt.Value = p.parseExpression()
	if err := p.expectStatementTerminator(); err != nil {
		return nil
	}
	return stmt
}

func (p *parser) parseFunction() *ast.FunctionDeclaration {
	stmt := &ast.FunctionDeclaration{}
	p.advanceToken() // consume token.FUNCTION
	ident, err := p.expect(token.IDENT)
	if err != nil {
		return nil
	}
	stmt.Name = &ast.Identifier{Value: ident.Literal}
	if _, err := p.expect(token.LPAREN); err != nil {
		return nil
	}
	if _, err := p.expect(token.RPAREN); err != nil {
		return nil
	}
	if _, err := p.expect(token.LBRACE); err != nil {
		return nil
	}
	stmt.Body = p.parseBody()
	if _, err := p.expect(token.RBRACE); err != nil {
		return nil
	}
	return stmt
}
