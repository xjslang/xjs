package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

type Builder struct {
	infixOperators   map[token.TokenType]infixOperator
	tokenizers       []func(l *lexer.Lexer, next func() token.Token) token.Token
	statementParsers []func(p *Parser, next func() (ast.Statement, error)) (ast.Statement, error)
}

func (b *Builder) Install(plugin func(p *Builder)) *Builder {
	plugin(b)
	return b
}

func (b *Builder) UseTokenizer(tokenizer func(l *lexer.Lexer, next func() token.Token) token.Token) {
	b.tokenizers = append(b.tokenizers, tokenizer)
}

// TODO: Builder.RegisterInfixOperator silently overwrites any existing entry for the same token type and does not validate inputs (e.g., negative precedence / nil fn), which can make plugin ordering bugs hard to detect and defer failures to a later panic during Build(). Consider detecting duplicates and validating arguments at registration time (and returning an error) to fail fast with clearer context.
func (b *Builder) RegisterInfixOperator(tt token.TokenType, precedence int, fn func(op token.Token, left, right ast.Expression) ast.Expression) {
	if b.infixOperators == nil {
		b.infixOperators = make(map[token.TokenType]infixOperator)
	}
	b.infixOperators[tt] = infixOperator{
		precedence: precedence,
		fn:         fn,
	}
}

func (b *Builder) UseStatementParser(parser func(p *Parser, next func() (ast.Statement, error)) (ast.Statement, error)) {
	b.statementParsers = append(b.statementParsers, parser)
}

func (b *Builder) Build(l *lexer.Lexer) *Parser {
	p := &Parser{}
	for _, tokenizer := range b.tokenizers {
		l.UseTokenizer(tokenizer)
	}
	for tt, infixOp := range b.infixOperators {
		p.RegisterInfixOperator(tt, infixOp.precedence, infixOp.fn)
	}
	for _, stmtParser := range b.statementParsers {
		p.UseStatementParser(stmtParser)
	}
	p.Init(l)
	return p
}
