package lexer

import (
	"github.com/xjslang/xjs/token"
)

type Builder struct {
	middlewares   []func(l *Lexer, next func() token.Token) token.Token
	dynamicTokens map[string]token.Type
	nextTokenID   token.Type
}

func NewBuilder() *Builder {
	return &Builder{
		middlewares:   []func(l *Lexer, next func() token.Token) token.Token{},
		dynamicTokens: make(map[string]token.Type),
		nextTokenID:   token.DYNAMIC_TOKENS_START,
	}
}

func (lb *Builder) Build(input string) *Lexer {
	return newWithOptions(input, lb.middlewares...)
}

func (lb *Builder) UseTokenReader(middleware func(l *Lexer, next func() token.Token) token.Token) *Builder {
	lb.middlewares = append(lb.middlewares, middleware)
	return lb
}

func (lb *Builder) RegisterTokenType(name string) token.Type {
	if tokenType, exists := lb.dynamicTokens[name]; exists {
		return tokenType
	}

	tokenType := lb.nextTokenID
	lb.nextTokenID++
	lb.dynamicTokens[name] = tokenType
	return tokenType
}
