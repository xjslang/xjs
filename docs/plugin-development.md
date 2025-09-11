# Plugin Development Guide

Learn how to create powerful extensions for XJS using the middleware system.

## Table of Contents

- [Plugin Architecture](#plugin-architecture)
- [Creating Your First Plugin](#creating-your-first-plugin)
- [Plugin Types](#plugin-types)
- [Best Practices](#best-practices)
- [Publishing Plugins](#publishing-plugins)
- [Advanced Patterns](#advanced-patterns)

## Plugin Architecture

XJS plugins work by registering middlewares that intercept and modify the parsing process. The parser processes middlewares in LIFO (Last-In, First-Out) order, allowing plugins to:

- **Override** default behavior
- **Transform** existing syntax
- **Add** new language features
- **Compose** with other plugins

```
Input → Lexer → Parser → [Plugin 3] → [Plugin 2] → [Plugin 1] → Default Parser → AST
```

## Creating Your First Plugin

Let's create a plugin that adds `const` keyword support.

### Step 1: Project Setup

```bash
mkdir xjs-const-plugin
cd xjs-const-plugin
go mod init github.com/yourname/xjs-const-plugin
go get github.com/xjslang/xjs@latest
```

### Step 2: Define the AST Node

```go
// const.go
package constplugin

import (
	"strings"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/token"
)

// ConstStatement represents a const declaration
type ConstStatement struct {
	Token token.Token     // the 'const' token
	Name  *ast.Identifier // variable name
	Value ast.Expression  // assigned value
}

// WriteTo implements the ast.Node interface
func (cs *ConstStatement) WriteTo(b *strings.Builder) {
	b.WriteString("const ")
	cs.Name.WriteTo(b)
	if cs.Value != nil {
		b.WriteString("=")
		cs.Value.WriteTo(b)
	}
}
```

### Step 3: Implement the Plugin

```go
// plugin.go
package constplugin

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// Plugin adds const keyword support
type Plugin struct{}

// New creates a new const plugin instance
func New() *Plugin {
	return &Plugin{}
}

// Apply registers the plugin with the parser
func (p *Plugin) Apply(parser *parser.Parser) {
	parser.UseStatementParser(p.parseConstStatement)
}

// parseConstStatement handles const declarations
func (p *Plugin) parseConstStatement(parser *parser.Parser, next func() ast.Statement) ast.Statement {
	if parser.CurrentToken.Type == token.IDENT && parser.CurrentToken.Literal == "const" {
		stmt := &ConstStatement{Token: parser.CurrentToken}
		
		// Move to identifier
		parser.NextToken()
		if parser.CurrentToken.Type != token.IDENT {
			return nil // Error: expected identifier
		}
		
		stmt.Name = &ast.Identifier{
			Token: parser.CurrentToken,
			Value: parser.CurrentToken.Literal,
		}
		
		// Expect assignment operator
		if !parser.ExpectToken(token.ASSIGN) {
			return nil
		}
		
		// Move to value expression
		parser.NextToken()
		stmt.Value = parser.ParseExpression()
		
		return stmt
	}
	
	return next() // Not our keyword, delegate to next middleware
}
```

### Step 4: Usage Example

```go
// example/main.go
package main

import (
	"fmt"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/yourname/xjs-const-plugin"
)

func main() {
	input := `const PI = 3.14159`
	
	l := lexer.New(input)
	p := parser.New(l)
	
	// Apply the plugin
	plugin := constplugin.New()
	plugin.Apply(p)
	
	program, err := p.ParseProgram()
	if err != nil {
		panic(err)
	}
	
	fmt.Println(program.String())
	// Output: const PI=3.14159
}
```

## Plugin Types

### 1. Statement Plugins

Add new statement types (declarations, control flow, etc.).

```go
type StatementPlugin struct{}

func (sp *StatementPlugin) Apply(p *parser.Parser) {
	p.UseStatementParser(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		// Check for your custom statement syntax
		if p.CurrentToken.Literal == "mystatement" {
			return sp.parseMyStatement(p)
		}
		return next()
	})
}
```

### 2. Expression Plugins

Add new expression types (operators, literals, etc.).

```go
type ExpressionPlugin struct{}

func (ep *ExpressionPlugin) Apply(p *parser.Parser) {
	p.UseExpressionParser(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Literal == "SPECIAL" {
			return &SpecialExpression{Token: p.CurrentToken}
		}
		return next()
	})
}
```

### 3. Operator Plugins

Add custom operators with proper precedence.

```go
type OperatorPlugin struct{}

func (op *OperatorPlugin) Apply(p *parser.Parser) {
	// Add prefix operator
	p.RegisterPrefixOperator("~", func(right func() ast.Expression) ast.Expression {
		return &BitwiseNotExpression{Right: right()}
	})
	
	// Add infix operator
	p.RegisterInfixOperator("**", parser.PRODUCT+1, 
		func(left ast.Expression, right func() ast.Expression) ast.Expression {
			return &PowerExpression{Left: left, Right: right()}
		})
}
```

### 4. Literal Plugins

Add custom literals and constants.

```go
type LiteralPlugin struct{}

func (lp *LiteralPlugin) Apply(p *parser.Parser) {
	// Add mathematical constants
	p.RegisterOperand("PI", func() ast.Expression {
		return &MathConstant{Name: "PI", Value: "Math.PI"}
	})
	
	p.RegisterOperand("E", func() ast.Expression {
		return &MathConstant{Name: "E", Value: "Math.E"}
	})
}
```

## Best Practices

### 1. Error Handling

Always provide meaningful error messages:

```go
func (p *Plugin) parseMyStatement(parser *parser.Parser, next func() ast.Statement) ast.Statement {
	if parser.CurrentToken.Literal == "mystmt" {
		stmt := &MyStatement{Token: parser.CurrentToken}
		
		parser.NextToken()
		if parser.CurrentToken.Type != token.IDENT {
			// Add error to parser
			parser.AddError(fmt.Sprintf("expected identifier, got %s", 
				parser.CurrentToken.Literal))
			return nil
		}
		
		// Continue parsing...
		return stmt
	}
	return next()
}
```

### 2. Token Preservation

Always preserve token information for debugging:

```go
type MyNode struct {
	Token token.Token // Always include the token
	Name  string
	Value ast.Expression
}
```

### 3. Graceful Fallback

Handle edge cases and provide fallbacks:

```go
func (p *Plugin) middleware(parser *parser.Parser, next func() ast.Statement) ast.Statement {
	if parser.CurrentToken.Type == token.EOF {
		return nil // Gracefully handle end of file
	}
	
	if parser.CurrentToken.Literal == "mykeyword" {
		// Validate syntax before proceeding
		if !p.isValidSyntax(parser) {
			return next() // Fall back to default parsing
		}
		return p.parseMyKeyword(parser)
	}
	
	return next()
}
```

### 4. Testing

Write comprehensive tests for your plugin:

```go
// plugin_test.go
package constplugin

import (
	"testing"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func TestConstDeclaration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`const x = 5`, `const x=5`},
		{`const name = "test"`, `const name="test"`},
		{`const arr = [1, 2, 3]`, `const arr=[1,2,3]`},
	}
	
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		
		plugin := New()
		plugin.Apply(p)
		
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("Parse error: %v", err)
		}
		
		if program.String() != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, program.String())
		}
	}
}

func TestConstErrors(t *testing.T) {
	invalidInputs := []string{
		`const`,           // Missing identifier
		`const 123`,       // Invalid identifier  
		`const x`,         // Missing assignment
		`const x =`,       // Missing value
	}
	
	for _, input := range invalidInputs {
		l := lexer.New(input)
		p := parser.New(l)
		
		plugin := New()
		plugin.Apply(p)
		
		_, err := p.ParseProgram()
		if err == nil && len(p.Errors()) == 0 {
			t.Errorf("Expected error for input: %s", input)
		}
	}
}
```

### 5. Documentation

Document your plugin thoroughly:

```go
// Package constplugin provides const keyword support for XJS
//
// This plugin adds support for const declarations which are transpiled
// to let declarations with a comment indicating immutability.
//
// Example usage:
//
//   plugin := constplugin.New()
//   plugin.Apply(parser)
//
//   // Now you can parse: const x = 5
//   // Which outputs: /* const */ let x=5
package constplugin
```

## Publishing Plugins

### 1. Repository Structure

```
xjs-myplugin/
├── README.md
├── go.mod
├── plugin.go
├── plugin_test.go
├── examples/
│   └── main.go
└── docs/
    └── api.md
```

### 2. Naming Convention

Follow the naming pattern: `xjs-<feature>-plugin` or `<feature>-parser`

Examples:
- `xjs-const-plugin`
- `jsx-parser`
- `typescript-parser`

### 3. README Template

```markdown
# XJS MyFeature Plugin

Adds [feature description] support to XJS.

## Installation

```bash
go get github.com/yourname/xjs-myfeature-plugin
```

## Usage

```go
import "github.com/yourname/xjs-myfeature-plugin"

plugin := myfeatureplugin.New()
plugin.Apply(parser)
```

## Examples

[Provide clear examples]

## API

[Document your API]

## Contributing

[Contribution guidelines]
```

### 4. Go Module

```go
module github.com/yourname/xjs-myplugin

go 1.21

require github.com/xjslang/xjs v1.0.0
```

## Advanced Patterns

### 1. Composable Plugins

Create plugins that work well together:

```go
type ComposablePlugin struct {
	priority int
	name     string
}

func (cp *ComposablePlugin) Priority() int {
	return cp.priority
}

func (cp *ComposablePlugin) Name() string {
	return cp.name
}

// Plugin manager
type PluginManager struct {
	plugins []ComposablePlugin
}

func (pm *PluginManager) Add(plugin ComposablePlugin) {
	pm.plugins = append(pm.plugins, plugin)
	// Sort by priority
	sort.Slice(pm.plugins, func(i, j int) bool {
		return pm.plugins[i].Priority() > pm.plugins[j].Priority()
	})
}
```

### 2. Configuration Support

Allow plugin configuration:

```go
type PluginConfig struct {
	StrictMode bool
	Transform  func(string) string
}

type ConfigurablePlugin struct {
	config PluginConfig
}

func NewWithConfig(config PluginConfig) *ConfigurablePlugin {
	return &ConfigurablePlugin{config: config}
}
```

### 3. Context-Aware Parsing

Use context to make parsing decisions:

```go
type Context struct {
	InsideFunction bool
	InsideLoop     bool
	Variables      map[string]bool
}

func (p *Plugin) parseWithContext(parser *parser.Parser, ctx *Context) ast.Statement {
	if ctx.InsideFunction && parser.CurrentToken.Literal == "yield" {
		return p.parseYieldStatement(parser)
	}
	
	if ctx.InsideLoop && parser.CurrentToken.Literal == "break" {
		return p.parseBreakStatement(parser)
	}
	
	return nil
}
```

### 4. Multi-pass Parsing

Some features require multiple parsing passes:

```go
type MultiPassPlugin struct {
	firstPassData map[string]interface{}
}

func (mp *MultiPassPlugin) FirstPass(program *ast.Program) {
	// Collect information in first pass
	mp.collectDeclarations(program)
}

func (mp *MultiPassPlugin) SecondPass(program *ast.Program) {
	// Use collected information to transform
	mp.resolveReferences(program)
}
```

### 5. Source Maps

Generate source maps for debugging:

```go
type SourceMapPlugin struct {
	mappings []SourceMapping
}

type SourceMapping struct {
	OriginalLine   int
	OriginalColumn int
	GeneratedLine  int
	GeneratedColumn int
}

func (smp *SourceMapPlugin) AddMapping(original, generated token.Position) {
	smp.mappings = append(smp.mappings, SourceMapping{
		OriginalLine:    original.Line,
		OriginalColumn:  original.Column,
		GeneratedLine:   generated.Line,
		GeneratedColumn: generated.Column,
	})
}
```

## Example: Complete JSX Plugin

Here's a simplified JSX plugin implementation:

```go
package jsxplugin

import (
	"strings"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

type JSXElement struct {
	Token      token.Token
	TagName    string
	Attributes map[string]ast.Expression
	Children   []ast.Expression
}

func (jsx *JSXElement) WriteTo(b *strings.Builder) {
	b.WriteString("React.createElement(\"")
	b.WriteString(jsx.TagName)
	b.WriteString("\"")
	
	if len(jsx.Attributes) > 0 {
		b.WriteString(",{")
		first := true
		for attr, value := range jsx.Attributes {
			if !first {
				b.WriteString(",")
			}
			b.WriteString(attr)
			b.WriteString(":")
			value.WriteTo(b)
			first = false
		}
		b.WriteString("}")
	} else {
		b.WriteString(",null")
	}
	
	for _, child := range jsx.Children {
		b.WriteString(",")
		child.WriteTo(b)
	}
	
	b.WriteString(")")
}

type Plugin struct{}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Apply(parser *parser.Parser) {
	parser.UseExpressionParser(p.parseJSX)
}

func (p *Plugin) parseJSX(parser *parser.Parser, next func() ast.Expression) ast.Expression {
	if parser.CurrentToken.Type == token.LT {
		// This might be JSX, check next token
		if parser.PeekToken().Type == token.IDENT {
			return p.parseJSXElement(parser)
		}
	}
	return next()
}

func (p *Plugin) parseJSXElement(parser *parser.Parser) ast.Expression {
	jsx := &JSXElement{
		Token:      parser.CurrentToken,
		Attributes: make(map[string]ast.Expression),
	}
	
	parser.NextToken() // skip <
	jsx.TagName = parser.CurrentToken.Literal
	
	// Parse attributes and children...
	// (Implementation details omitted for brevity)
	
	return jsx
}
```

---

## Community and Support

- **GitHub Discussions**: Share ideas and get help
- **Plugin Registry**: List your plugin in the official registry
- **Examples Repository**: Submit examples of your plugin in action

Ready to create your first plugin? Start with the [const plugin example](examples/const-plugin/) and build from there!