package lexer

import (
	"github.com/xjslang/xjs/token"
)

// NewBuilder creates a new Builder instance for constructing customized lexers.
func NewBuilder() *Builder {
	return &Builder{
		interceptors:  []Interceptor{},
		dynamicTokens: make(map[string]token.Type),
		nextTokenID:   token.DYNAMIC_TOKENS_START,
	}
}

// Build creates a new Lexer instance with the configured interceptors and dynamic tokens.
func (lb *Builder) Build(input string) *Lexer {
	return newWithOptions(input, lb.interceptors...)
}

// UseTokenInterceptor adds a middleware interceptor to the lexer pipeline.
func (lb *Builder) UseTokenInterceptor(interceptor Interceptor) *Builder {
	lb.interceptors = append(lb.interceptors, interceptor)
	return lb
}

// RegisterTokenType creates and registers a new dynamic token type.
func (lb *Builder) RegisterTokenType(name string) token.Type {
	if tokenType, exists := lb.dynamicTokens[name]; exists {
		return tokenType
	}

	tokenType := lb.nextTokenID
	lb.nextTokenID++
	lb.dynamicTokens[name] = tokenType
	return tokenType
}
