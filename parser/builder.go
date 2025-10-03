package parser

import (
	"fmt"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/token"
)

func NewBuilder(lb *lexer.Builder) *Builder {
	// initializes maps with existing operators from precedences
	registeredInfixOps := make(map[token.Type]bool)
	for tokenType := range precedences {
		registeredInfixOps[tokenType] = true
	}

	// initializes map for prefix operators (from parser's default prefix operators)
	registeredPrefixOps := map[token.Type]bool{
		token.IDENT:      true,
		token.INT:        true,
		token.FLOAT:      true,
		token.STRING:     true,
		token.RAW_STRING: true,
		token.TRUE:       true,
		token.FALSE:      true,
		token.NULL:       true,
		token.NOT:        true,
		token.MINUS:      true,
		token.INCREMENT:  true,
		token.DECREMENT:  true,
		token.LPAREN:     true,
		token.LBRACKET:   true,
		token.LBRACE:     true,
		token.FUNCTION:   true,
	}

	return &Builder{
		LexerBuilder:        lb,
		stmtInterceptors:    []Interceptor[ast.Statement]{},
		expInterceptors:     []Interceptor[ast.Expression]{},
		prefixOperators:     []prefixOperator{},
		infixOperators:      []infixOperator{},
		registeredPrefixOps: registeredPrefixOps,
		registeredInfixOps:  registeredInfixOps,
	}
}

// Install installs a plugin, allowing to enhance the language by modifying or adding new operators, statements, and expressions.
func (pb *Builder) Install(plugin func(*Builder)) *Builder {
	plugin(pb)
	return pb
}

// UseStatementInterceptor intercepts the parsing flow for statements, allowing you to modify or add custom statements.
func (pb *Builder) UseStatementInterceptor(interceptor Interceptor[ast.Statement]) *Builder {
	pb.stmtInterceptors = append(pb.stmtInterceptors, interceptor)
	return pb
}

// UseExpressionInterceptor intercepts the parsing flow for expressions, allowing you to modify or add custom expressions.
func (pb *Builder) UseExpressionInterceptor(interceptor Interceptor[ast.Expression]) *Builder {
	pb.expInterceptors = append(pb.expInterceptors, interceptor)
	return pb
}

// RegisterPrefixOperator allows registering new prefix operators (for example, `typeof`).
func (pb *Builder) RegisterPrefixOperator(tokenType token.Type, createExpr func(tok token.Token, right func() ast.Expression) ast.Expression) error {
	if pb.registeredPrefixOps[tokenType] {
		return fmt.Errorf("duplicate prefix operator: %s", tokenType)
	}
	pb.prefixOperators = append(pb.prefixOperators, prefixOperator{
		tokenType:  tokenType,
		createExpr: createExpr,
	})
	pb.registeredPrefixOps[tokenType] = true
	return nil
}

// RegisterInfixOperator allows incorporating new infix operators (for example, `^` for power).
func (pb *Builder) RegisterInfixOperator(tokenType token.Type, precedence int, createExpr func(tok token.Token, left ast.Expression, right func() ast.Expression) ast.Expression) error {
	if pb.registeredInfixOps[tokenType] {
		return fmt.Errorf("duplicate infix operator: %s", tokenType)
	}
	pb.infixOperators = append(pb.infixOperators, infixOperator{
		tokenType:  tokenType,
		precedence: precedence,
		createExpr: createExpr,
	})
	pb.registeredInfixOps[tokenType] = true
	return nil
}

// WithTolerantMode enables tolerant parsing mode, which continues parsing even on syntax errors.
// This is useful for language servers, formatters, and analysis tools that need to work with
// incomplete or invalid code. In tolerant mode, the parser will not stop on missing semicolons
// or other recoverable syntax errors.
//
// Example:
//
//	parser := parser.NewBuilder(lexer.NewBuilder()).
//	    WithTolerantMode(true).
//	    Build(code)
func (pb *Builder) WithTolerantMode(enabled bool) *Builder {
	pb.tolerantMode = enabled
	return pb
}

// Build creates a new instance of the parser.
func (pb *Builder) Build(input string) *Parser {
	l := pb.LexerBuilder.Build(input)
	return newWithOptions(l, parserOptions{
		stmtInterceptors: pb.stmtInterceptors,
		expInterceptors:  pb.expInterceptors,
		prefixOperators:  pb.prefixOperators,
		infixOperators:   pb.infixOperators,
		tolerantMode:     pb.tolerantMode,
	})
}
