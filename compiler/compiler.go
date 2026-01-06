package compiler

import (
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/sourcemap"
)

type CompileResult struct {
	Code      string
	SourceMap *sourcemap.SourceMap
}

// PrettyPrintOption configures formatting behavior.
// We use the functional options pattern here to provide sensible defaults
// while allowing extensibility as formatting options grow.
type PrettyPrintOptions struct {
	// IndentString is the string used for each indentation level
	// Common values: "  " (2 spaces), "    " (4 spaces), "\t" (tab)
	IndentString string
	// WriteSemicolons controls whether to write optional semicolons
	// true (default): write semicolons after statements
	// false: omit optional semicolons (required semicolons like in for loops are always written)
	WriteSemicolons bool
}

// PrettyPrintOption is a function that configures PrettyPrintOptions
type PrettyPrintOption func(*PrettyPrintOptions)

// WithSpaces creates an option with the specified number of spaces for indentation
func WithSpaces(count int) PrettyPrintOption {
	return func(opts *PrettyPrintOptions) {
		opts.IndentString = strings.Repeat(" ", count)
	}
}

// WithTabs creates an option with tabs for indentation
func WithTabs() PrettyPrintOption {
	return func(opts *PrettyPrintOptions) {
		opts.IndentString = "\t"
	}
}

func WithSemi(value bool) PrettyPrintOption {
	return func(opts *PrettyPrintOptions) {
		opts.WriteSemicolons = value
	}
}

type Compiler struct {
	generateSourceMap  bool
	prettyPrint        bool
	prettyPrintOptions PrettyPrintOptions
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) WithPrettyPrint(opts ...PrettyPrintOption) *Compiler {
	c.prettyPrint = true
	// Start with default options
	c.prettyPrintOptions = PrettyPrintOptions{
		IndentString:    "  ", // 2 spaces by default
		WriteSemicolons: true, // write semicolons by default
	}
	// Apply all provided options in order
	for _, opt := range opts {
		opt(&c.prettyPrintOptions)
	}
	return c
}

func (c *Compiler) WithSourceMap() *Compiler {
	c.generateSourceMap = true
	return c
}

func (c *Compiler) Compile(program *ast.Program) CompileResult {
	w := ast.CodeWriter{
		Builder:         strings.Builder{},
		PrettyPrint:     c.prettyPrint,
		IndentString:    c.prettyPrintOptions.IndentString,
		WriteSemicolons: c.prettyPrintOptions.WriteSemicolons,
	}
	if c.generateSourceMap {
		w.Mapper = sourcemap.New()
	}
	program.WriteTo(&w)

	code := w.String()
	var sm *sourcemap.SourceMap
	if c.generateSourceMap {
		sm = w.Mapper.SourceMap()
	}

	return CompileResult{Code: code, SourceMap: sm}
}
