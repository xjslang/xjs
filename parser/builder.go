package parser

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

type Builder struct {
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

func (b *Builder) UseStatementParser(parser func(p *Parser, next func() (ast.Statement, error)) (ast.Statement, error)) {
	b.statementParsers = append(b.statementParsers, parser)
}

func (b *Builder) Build(l *lexer.Lexer) *Parser {
	p := &Parser{}
	for _, tokenizer := range b.tokenizers {
		l.UseTokenizer(tokenizer)
	}
	for _, stmtParser := range b.statementParsers {
		p.UseStatementParser(stmtParser)
	}
	p.Init(l)
	return p
}
