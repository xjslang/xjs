package parser

import "slices"

type ContextType int

const (
	GlobalContext ContextType = iota
	FunctionContext
	BlockContext
)

func (p *XJSParser) PushContext(ctx ContextType) {
	p.contextStack = append(p.contextStack, ctx)
}

func (p *XJSParser) PopContext() {
	if len(p.contextStack) > 0 {
		p.contextStack = p.contextStack[:len(p.contextStack)-1]
	}
}

func (p *XJSParser) CurrentContext() ContextType {
	if len(p.contextStack) == 0 {
		return GlobalContext
	}
	return p.contextStack[len(p.contextStack)-1]
}

func (p *XJSParser) IsInFunction() bool {
	return slices.Contains(p.contextStack, FunctionContext)
}
