# XJS (eXtensible JavaScript Parser) - Copilot Instructions

## Project Overview

XJS is a highly customizable JavaScript parser written in Go. It's designed around the principle of **minimalism and sufficiency**, providing only essential JavaScript features while offering a powerful extension system through middleware patterns.

## Core Philosophy

- **Minimalist approach**: Only includes necessary and sufficient language constructs
- **Extensible architecture**: Everything can be customized via middlewares and plugins
- **No feature bloat**: Deliberately excludes redundant JavaScript features
- **Plugin-first design**: Missing features can be added through plugins

### Excluded Features (by design)
- No classes (use functions)
- No arrow functions (use regular functions)
- No `const/var` (only `let`)
- No `try/catch` (use alternative error handling)
- No redundant syntactic sugar

## Project Structure

```
xjs/
├── ast/           # Abstract Syntax Tree definitions
├── lexer/         # Tokenization and lexical analysis
├── parser/        # Syntax parsing and AST building
├── token/         # Token types and structures
├── internal/      # Internal utilities (debug, etc.)
├── test/          # Integration tests and examples
│   ├── integration/
│   └── testdata/
└── magefile.go    # Build automation
```

## Key Components

### 1. Lexer (`lexer/lexer.go`)
- Tokenizes JavaScript source code
- Tracks line and column positions
- Handles strings, numbers, identifiers, operators, and keywords

### 2. Parser (`parser/parser.go`)
- Builds AST from tokens
- Uses Pratt parsing for expression precedence
- Supports middleware pattern for extensibility
- Operator precedence levels: LOWEST → ASSIGNMENT → LOGICAL_OR → LOGICAL_AND → EQUALITY → COMPARISON → SUM → PRODUCT → UNARY → POSTFIX → CALL → MEMBER

### 3. AST (`ast/ast.go`)
- Defines all language constructs as nodes
- Each node implements `WriteTo(b *strings.Builder)` for code generation
- Two main interfaces: `Statement` and `Expression`

### 4. Token (`token/token.go`)
- Defines all token types used by lexer/parser
- Includes operators, keywords, literals, and delimiters

### 5. Error Handling (`parser/errors.go`)
- Structured error reporting with position information
- JSON serialization support for errors
- Multiple error collection and reporting

## Extension Mechanisms

### 1. UseStatementParser
Add custom statement parsing logic:
```go
p.UseStatementParser(func(p *Parser, next func() Statement) Statement {
    // Custom statement parsing logic
    return next() // fallback to default
})
```

### 2. UseExpressionParser
Add custom expression parsing logic:
```go
p.UseExpressionParser(func(p *Parser, next func() Expression) Expression {
    // Custom expression parsing logic
    return next() // fallback to default
})
```

### 3. UsePrefixOperator
Add custom prefix operators (like `typeof`, `!`, `-`):
```go
p.UsePrefixOperator("typeof", func(right func() Expression) Expression {
    // Return custom expression node
})
```

### 4. UseInfixOperator
Add custom infix operators with precedence:
```go
p.UseInfixOperator("^", PRODUCT+1, func(left Expression, right func() Expression) Expression {
    // Return custom expression node
})
```

### 5. UseOperand
Add custom operands/literals:
```go
p.UseOperand("PI", func() Expression {
    // Return custom expression node
})
```

## Development Guidelines

### Code Style
- Use Go conventions (gofmt, golint)
- Package comments should explain the package purpose
- Public functions need documentation
- Test coverage is important

### Adding New Features
1. **Consider if it belongs in core**: Most features should be plugins
2. **Write tests first**: Use test-driven development
3. **Update AST**: Add new node types to `ast/ast.go`
4. **Implement WriteTo**: Ensure code generation works
5. **Add parsing logic**: Update lexer/parser as needed
6. **Add integration tests**: Test end-to-end functionality

### Error Handling
- Use structured errors with position information
- Collect multiple errors when possible
- Provide helpful error messages
- Support JSON error output for tooling

### Testing
- Unit tests for individual components
- Integration tests in `test/integration/`
- Test data in `test/testdata/` with `.js` and `.output` files:
  - `.js` files contain the source code to be parsed
  - `.output` files contain the **execution result** of running the parsed/transpiled code (NOT the transpiled code itself)
- Example-based tests using Go's Example functions

## Common Patterns

### Creating New AST Nodes
```go
type MyNode struct {
    Token token.Token
    Field Expression // or other fields
}

func (mn *MyNode) WriteTo(b *strings.Builder) {
    // Implement code generation
    b.WriteString("my_output")
}
```

### Parser Middleware Pattern
```go
// Middleware executed in LIFO order (Last-In, First-Out)
p.UseStatementParser(middleware1) // executed last
p.UseStatementParser(middleware2) // executed first
```

### Operator Precedence
Always consider precedence when adding new operators:
- Higher number = higher precedence
- Use existing constants as reference points
- Test with complex expressions

## Dependencies
- `github.com/magefile/mage`: Build automation
- `github.com/dop251/goja`: JavaScript engine (for testing/validation)
- Standard Go library only for core functionality

## Build and Test
- Use `mage` for build automation
- Run tests with `go test ./...`
- Integration tests validate end-to-end behavior
- Test data files provide expected outputs

## Plugin Development
When creating plugins (separate repositories):
- Follow the middleware patterns
- Implement proper AST nodes with WriteTo methods
- Provide comprehensive tests
- Document usage examples
- Consider operator precedence carefully

## Key Files to Understand
- `parser/parser.go`: Core parsing logic
- `parser/parser_functions.go`: Expression and statement parsing
- `parser/parser_middlewares.go`: Middleware system
- `ast/ast.go`: All AST node definitions
- `lexer/lexer.go`: Tokenization logic
- `token/token.go`: Token type definitions

## Tips for Contributors
1. Start with understanding the existing AST nodes
2. Look at integration tests for usage patterns
3. The middleware system is key to extensibility
4. Always implement WriteTo for code generation
5. Consider both parsing and code generation phases
6. Test with realistic JavaScript-like code examples