package parser

import "slices"

// ContextType represents different parsing contexts that help track the current
// scope during parsing for proper semantic analysis.
type ContextType int

const (
	// GlobalContext represents the top-level global scope
	GlobalContext ContextType = iota
	// FunctionContext represents parsing inside a function body
	FunctionContext
	// BlockContext represents parsing inside a block statement
	BlockContext
)

// PushContext adds a new context to the top of the context stack.
func (p *Parser) PushContext(ctx ContextType) {
	p.contextStack = append(p.contextStack, ctx)
}

// PopContext removes the top context from the context stack.
func (p *Parser) PopContext() {
	if len(p.contextStack) > 0 {
		p.contextStack = p.contextStack[:len(p.contextStack)-1]
	}
}

// CurrentContext returns the context at the top of the stack.
func (p *Parser) CurrentContext() ContextType {
	if len(p.contextStack) == 0 {
		return GlobalContext
	}
	return p.contextStack[len(p.contextStack)-1]
}

// IsInFunction checks if there's a FunctionContext anywhere in the context stack.
func (p *Parser) IsInFunction() bool {
	return slices.Contains(p.contextStack, FunctionContext)
}
