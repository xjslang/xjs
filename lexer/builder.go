package lexer

import (
	"github.com/xjslang/xjs/token"
)

func NewBuilder() *Builder {
	return &Builder{
		interceptors:  []Interceptor{},
		dynamicTokens: make(map[string]token.Type),
		nextTokenID:   token.DYNAMIC_TOKENS_START,
	}
}

func (lb *Builder) Build(input string) *Lexer {
	return newWithOptions(input, lb.interceptors...)
}

func (lb *Builder) UseInterceptor(interceptor Interceptor) *Builder {
	lb.interceptors = append(lb.interceptors, interceptor)
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
