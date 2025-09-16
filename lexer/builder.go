package lexer

import (
	"github.com/xjslang/xjs/token"
)

func NewBuilder() *Builder {
	return &Builder{
		readers:       []Reader{},
		dynamicTokens: make(map[string]token.Type),
		nextTokenID:   token.DYNAMIC_TOKENS_START,
	}
}

func (lb *Builder) Build(input string) *Lexer {
	return newWithOptions(input, lb.readers...)
}

func (lb *Builder) UseTokenReader(reader Reader) *Builder {
	lb.readers = append(lb.readers, reader)
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
