package parser

import "slices"

type ContextType int

const (
	GlobalContext ContextType = iota
	FunctionContext
	BlockContext
)

func (p *Parser) PushContext(ctx ContextType) {
	p.contextStack = append(p.contextStack, ctx)
}

func (p *Parser) PopContext() {
	if len(p.contextStack) > 0 {
		p.contextStack = p.contextStack[:len(p.contextStack)-1]
	}
}

func (p *Parser) CurrentContext() ContextType {
	if len(p.contextStack) == 0 {
		return GlobalContext
	}
	return p.contextStack[len(p.contextStack)-1]
}

func (p *Parser) IsInFunction() bool {
	return slices.Contains(p.contextStack, FunctionContext)
}
