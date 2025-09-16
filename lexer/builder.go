package lexer

import (
	"github.com/xjslang/xjs/token"
)

type XJSBuilder struct {
	readers       []func(l Lexer, next func() token.Token) token.Token
	dynamicTokens map[string]token.Type
	nextTokenID   token.Type
}

func NewBuilder() *XJSBuilder {
	return &XJSBuilder{
		readers:       []func(l Lexer, next func() token.Token) token.Token{},
		dynamicTokens: make(map[string]token.Type),
		nextTokenID:   token.DYNAMIC_TOKENS_START,
	}
}

func (lb *XJSBuilder) Build(input string) *XJSLexer {
	return newWithOptions(input, lb.readers...)
}

func (lb *XJSBuilder) UseTokenReader(reader func(l Lexer, next func() token.Token) token.Token) *XJSBuilder {
	lb.readers = append(lb.readers, reader)
	return lb
}

func (lb *XJSBuilder) RegisterTokenType(name string) token.Type {
	if tokenType, exists := lb.dynamicTokens[name]; exists {
		return tokenType
	}

	tokenType := lb.nextTokenID
	lb.nextTokenID++
	lb.dynamicTokens[name] = tokenType
	return tokenType
}
