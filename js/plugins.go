package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

var (
	letType  = token.RegisterType("let")
	funcType = token.RegisterType("function")
)

func LetPlugin(b *parser.Builder) {
	b.UseTokenizer(func(l *lexer.Lexer, next func() token.Token) token.Token {
		tok := next()
		if lit := letType.String(); tok.Type == token.IDENT && tok.Literal == lit {
			return token.Token{Type: letType, Literal: lit}
		}
		return tok
	})
	b.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == letType {
			stmt, err := ParseLetStatement(p)
			if err != nil {
				return nil, err
			}
			return stmt, nil
		}
		return next()
	})
}

func FunctionPlugin(b *parser.Builder) {
	b.UseTokenizer(func(l *lexer.Lexer, next func() token.Token) token.Token {
		tok := next()
		if lit := funcType.String(); tok.Type == token.IDENT && tok.Literal == lit {
			return token.Token{Type: funcType, Literal: lit}
		}
		return tok
	})
	b.UseStatementParser(func(p *parser.Parser, next func() (ast.Statement, error)) (ast.Statement, error) {
		if p.CurrentToken.Type == funcType {
			stmt, err := ParseFunctionDeclaration(p)
			if err != nil {
				return nil, err
			}
			return stmt, nil
		}
		return next()
	})
}
